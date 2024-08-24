package pico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/http"
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
}

func (p *Pico) receive() {
	header := make([]byte, HeaderSize)
	buf := make([]byte, BufferSize)

	//bufio.NewReaderSize(p.conn, BufferSize)
	reader := bufio.NewReader(p.conn)

	for {
		n, e := reader.Read(header)
		if e != nil {
			break
		}
		if n < HeaderSize {
			//break
			continue //TODO继续接受
		}

		if bytes.Compare(header[:MagicSize], []byte(Magic)) != 0 {
			continue
		}

		//解析
		id := binary.BigEndian.Uint16(header[4:])
		typ := header[6] >> 4
		encoding := header[6] & 0x0f

		length := int(header[7])<<16 + int(header[8])<<8 + int(header[9])
		if length > 0 {
			var b []byte
			if length > BufferSize {
				b = make([]byte, length)
			} else {
				b = buf
			}

			//_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * 30))
			n, err := io.ReadAtLeast(reader, b, length)
			if err != nil {
				break
			}
			if n != length {
				//长度不够，废包
				break
			}
		}

		switch typ {

		}
	}
}

func (p *Pico) AttachHandler(h http.Handler) {
	p.handler = h
}

func (p *Pico) handleRequest(buf []byte) {
	if p.handler == nil {
		return
	}

	//1 解析request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	if err != nil {
		return
	}

	//构建response，接收响应
	rw := newHttpResponseWriter()

	//2 执行请求
	//req.Header.Set("token", "inline") //使用内置token，免验证
	p.handler.ServeHTTP(rw, req)

	//3 回传 response
	rw.Bytes()
}
