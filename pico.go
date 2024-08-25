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
	"sync/atomic"
	"time"
)

const (
	Magic      = "pico"
	MagicSize  = len(Magic)
	HeaderSize = 10
	BufferSize = 1024
)

type Pico struct {
	handler http.Handler
	conn    net.Conn
	id      atomic.Uint32

	reader     *bufio.Reader
	writer     *bufio.Writer
	writerLock sync.Mutex

	streams  Map[uint16, stream]
	requests Map[uint16, pending]
	pings    Map[uint16, pending]

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
	return uint16(p.id.Add(1))
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
		} else {
			//todo 内存复制问题
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

	_, err := p.writer.Write(header)
	if err != nil {
		return err
	}

	_, err = p.writer.Write(pack.Payload)
	if err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *Pico) Request(req *http.Request) (*http.Response, error) {
	buf := bytes.NewBuffer(nil)
	err := req.Write(buf)
	if err != nil {
		return nil, err
	}

	id := p.getId()

	err = p.Send(&Pack{Id: id, Type: REQUEST, Payload: buf.Bytes()})
	if err != nil {
		return nil, err
	}

	r := newPending()
	p.requests.Store(id, r)
	defer p.requests.Delete(id)

	select {
	case <-time.After(time.Minute):
		return nil, errors.New("ping timeout")
	case pack := <-r.c:
		if pack == nil {
			return nil, errors.New("invalid pending")
		}
		return http.ReadResponse(bufio.NewReader(bytes.NewReader(pack.Payload)), req)
	}
}

func (p *Pico) Ping() (ms int, err error) {
	id := p.getId()
	err = p.Send(&Pack{Id: id, Type: PING})
	if err != nil {
		return
	}

	start := time.Now().UnixMilli()

	r := newPending()
	p.pings.Store(id, r)
	defer p.pings.Delete(id)

	select {
	case <-time.After(time.Minute):
		err = errors.New("ping timeout")
		return
	case pack := <-r.c:
		if pack == nil {
			err = errors.New("invalid pending")
			return
		}
		return int(time.Now().UnixMilli() - start), nil
	}
}

func (p *Pico) Stream() (rw io.ReadWriteCloser, id uint16) {
	id = p.getId()
	stream := newStream(p, id)
	p.streams.Store(id, stream)
	return stream, id
}

func (p *Pico) handlePing(pack *Pack) {
	pack.Type = PONG
	err := p.Send(pack)
	if err != nil {
		//todo log
	}
}

func (p *Pico) handlePong(pack *Pack) {
	pending := p.pings.Load(pack.Id)
	if pending != nil {
		pending.c <- pack
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

func (p *Pico) handleResponse(pack *Pack) {
	pending := p.requests.Load(pack.Id)
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

func (p *Pico) handleStreamEnd(pack *Pack) {
	stream := p.streams.Load(pack.Id)
	if stream != nil {
		stream.put(pack)
	}
}
