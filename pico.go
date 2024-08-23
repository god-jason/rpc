package pico

import (
	"bufio"
	"bytes"
	"net/http"
)

type Pico struct {
	handler http.Handler
}

func (p *Pico) AttachHandler(h http.Handler) {
	p.handler = h
}

func (p *Pico) Connect(addr string) error {

}

func (p *Pico) Serve(port int) error {

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
