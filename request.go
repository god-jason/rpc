package pico

import (
	"bufio"
	"bytes"
	"net/http"
)

func Request(req *http.Request) (resp *http.Response, err error) {

	buffer := bytes.NewBuffer(nil)
	err = req.Write(buffer)
	if err != nil {
		return
	}

	//2 发送 request
	pico.Send(buffer.Bytes())

	//3 接收 response
	pico.Receive(buffer)

	var buf []byte

	//4 解析 response
	resp, err = http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), req)
	return
}

func AttachHandler(handler http.Handler) {
	pico.httpHandler = handler
}

func handleRequest(buf []byte) {

	//1 解析request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	if err != nil {
		return
	}

	//构建response，接收响应
	rw := newHttpResponseWriter()

	//2 执行请求
	//req.Header.Set("token", "inline") //使用内置token，免验证
	handler.ServeHTTP(rw, req)

	//3 回传 response
	rw.Bytes()

}
