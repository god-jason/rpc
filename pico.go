package pico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	Magic      = "pico"
	MagicSize  = len(Magic)
	HeaderSize = 10
	BufferSize = 1024
)

const (
	DISCONNECT uint8 = iota
	CONNECT
	CONNECT_ACK
	PING
	PONG
	REQUEST
	RESPONSE
	STREAM
	STREAM_END
	PUBLISH
	PUBLISH_ACK
	SUBSCRIBE
	SUBSCRIBE_ACK
	UNSUBSCRIBE
	UNSUBSCRIBE_ACK
	MESSAGE
)

const (
	BINARY uint8 = iota
	JSON
	XML
	YAML
	MSGPACK
)

type Pico struct {
	handler http.Handler
	conn    net.Conn
	id      uint16
	idLock  sync.Mutex

	reader     *bufio.Reader
	writer     *bufio.Writer
	writerLock sync.Mutex

	streams Map[uint16, Stream]

	pending Map[uint16, pending]

	header []byte
	buf    []byte
}

func (p *Pico) init() {
	p.header = make([]byte, HeaderSize)
	p.buf = make([]byte, BufferSize)
}

func (p *Pico) attach(conn net.Conn) {
	p.conn = conn
	p.reader = bufio.NewReader(conn)
	p.writer = bufio.NewWriter(conn)
}

func (p *Pico) getId() uint16 {
	p.idLock.Lock()
	defer p.idLock.Unlock()

	//自增
	p.id++

	//避免与streams重复
	for p.streams.Load(p.id) != nil {
		p.id++
	}

	//避免与pending重复，一般不会发生那么长时间的等待
	for p.pending.Load(p.id) != nil {
		p.id++
	}

	return p.id
}

func (p *Pico) AttachHandler(h http.Handler) {
	p.handler = h
}

func (p *Pico) readPack() (*Pack, error) {
	n, e := p.reader.Read(p.header)
	if e != nil {
		return nil, e
	}
	if n < HeaderSize {
		//break
		//continue //TODO继续接受
		//c.UnreadByte()
		return nil, errors.New("invalid header size")
	}

	if bytes.Compare(p.header[:MagicSize], []byte(Magic)) != 0 {
		return nil, errors.New("invalid header magic")
	}

	//解析
	pack := &Pack{
		Id:       binary.BigEndian.Uint16(p.header[4:]),
		Type:     p.header[6] >> 4,
		Encoding: p.header[6] & 0x0f,
	}

	length := int(p.header[7])<<16 + int(p.header[8])<<8 + int(p.header[9])
	if length > 0 {
		var b []byte
		if length > BufferSize {
			b = make([]byte, length)
		} else if pack.Type == STREAM || pack.Type == STREAM_END {
			//流数据，复制内存
			b = make([]byte, length)
		} else {
			//复用内存，可能会有问题
			b = p.buf
		}

		//_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		n, err := io.ReadAtLeast(p.reader, b, length)
		if err != nil {
			return nil, err
		}
		if n != length {
			//长度不够，废包
			return nil, errors.New("invalid data length")
		}

		pack.Payload = b[:n]
	}
	return pack, nil
}

func (p *Pico) Send(pack *Pack) error {
	//发送前，编码
	err := pack.Encode()
	if err != nil {
		return err
	}

	p.writerLock.Lock()
	defer p.writerLock.Unlock()

	//包头
	header := make([]byte, HeaderSize)
	copy(header[:MagicSize], []byte(Magic))
	binary.BigEndian.PutUint16(header[4:], pack.Id)
	header[6] = pack.Type<<4 + pack.Encoding
	length := len(pack.Payload)
	header[7] = byte(length >> 16)
	header[8] = byte(length >> 8)
	header[9] = byte(length)

	_, err = p.writer.Write(header)
	if err != nil {
		return err
	}

	_, err = p.writer.Write(pack.Payload)
	if err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *Pico) Ask(req *Pack) (resp *Pack, err error) {
	//自增序号
	req.Id = p.getId()

	//发送
	err = p.Send(req)
	if err != nil {
		return
	}

	r := newPending()
	p.pending.Store(req.Id, r)
	defer p.pending.Delete(req.Id)

	//等待结果
	select {
	case <-time.After(time.Minute): //TODO配置化
		err = errors.New("ping timeout")
		return
	case pack := <-r.c:
		if pack == nil {
			err = errors.New("invalid pending")
			return
		}
		return pack, nil
	}
}

func (p *Pico) Request(req *http.Request) (*http.Response, error) {
	buf := bytes.NewBuffer(nil)
	err := req.Write(buf)
	if err != nil {
		return nil, err
	}

	pack, err := p.Ask(&Pack{Type: REQUEST, Payload: buf.Bytes()})
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(bytes.NewReader(pack.Payload)), req)
}

func (p *Pico) Ping() (ms int, err error) {
	start := time.Now().UnixMilli()

	_, err = p.Ask(&Pack{Type: PING})
	if err != nil {
		return
	}

	return int(time.Now().UnixMilli() - start), nil
}

func (p *Pico) Stream() (stream *Stream, id uint16) {
	id = p.getId()
	stream = newStream(p, id)
	p.streams.Store(id, stream)
	return
}

func (p *Pico) handlePing(pack *Pack) {
	pack.Type = PONG
	err := p.Send(pack)
	if err != nil {
		//todo log
	}
}
func (p *Pico) handleRequest(pack *Pack) {
	if p.handler == nil {
		return
	}

	//1 解析request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(pack.Payload)))
	if err != nil {
		return
	}

	//构建response，接收响应
	rw := newHttpResponseWriter()

	//2 执行请求
	//req.Header.Set("token", "inline") //使用内置token，免验证
	p.handler.ServeHTTP(rw, req)

	//3 回传 response
	pack.Type = RESPONSE
	pack.Payload = rw.Bytes()

	err = p.Send(pack)
	if err != nil {

	}
}

func (p *Pico) handleAsk(pack *Pack) {
	pending := p.pending.Load(pack.Id)
	if pending != nil {
		pending.c <- pack
	}
}

func (p *Pico) handleStream(pack *Pack) {
	stream := p.streams.Load(pack.Id)
	if stream != nil {
		stream.put(pack)
	}
}
