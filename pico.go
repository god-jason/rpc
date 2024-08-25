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
	id      uint16

	reader     *bufio.Reader
	writer     *bufio.Writer
	writerLock sync.Mutex

	streams sync.Map

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
		//reader.UnreadByte()
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
