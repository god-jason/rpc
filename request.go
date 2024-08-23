package pico

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/textproto"
)

func Request() {

	var conn net.Conn

	//1 构建 request
	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	buf := req.Write(conn) //长度还不知道

	//2 发送 request
	pico.Send(buf)

	//3 接收 response
	pico.Receive(buf)

	//4 解析 response
	resp, err := http.ReadResponse(conn)

}

func handleRequest(buf []byte) {

	//1 解析request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))

	//构建response，接收响应
	rw := &ResponseWriter{}

	//2 执行请求
	app := gin.New()
	app.Run("")
	req.Header.Set("token", "inline") //使用内置token，免验证
	app.ServeHTTP(rw, req)

	//3 回传 response
	conn.Write(rw.Bytes())

}

type ResponseWriter struct {
	writer textproto.Writer
	header http.Header
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {

}

func (rw *ResponseWriter) WriteHeader(statusCode int) {

}

func (rw *ResponseWriter) Bytes() []byte {

}
