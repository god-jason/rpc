package pico

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
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

func handleRequest(buf []byte) {

	//1 解析request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	if err != nil {
		return
	}

	//构建response，接收响应
	rw := NewHttpResponseWriter()

	//2 执行请求
	app := gin.New()
	req.Header.Set("token", "inline") //使用内置token，免验证
	app.ServeHTTP(rw, req)

	//3 回传 response
	conn.Write(rw.Bytes())

}
