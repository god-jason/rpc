package pico

import (
	"bufio"
	"bytes"
	"net/http"
)

type Client struct {
	Pico
}

func (c *Client) Disconnect(reason string) {
	//todo 发送 disconnect
	_ = c.conn.Close()
}

func (c *Client) receive() {
	for {
		pack, err := c.readPack()
		if err != nil {
			break
		}
		c.handle(pack)
	}
}

func (c *Client) AttachHandler(h http.Handler) {
	c.handler = h
}

func (c *Client) handle(pack *Pack) {
	switch pack.Type {
	case CONNECT_ACK:
		c.handleConnectAck(pack)
	case PING:
		c.handlePing(pack)
	case PONG:
		c.handlePong(pack)
	case REQUEST:
		c.handleRequest(pack)
	case RESPONSE:
		c.handleResponse(pack)
	case STREAM:
		c.handleStream(pack)
	case STREAM_END:
		c.handleStreamEnd(pack)
	case PUBLISH:
		c.handlePublish(pack)
	case PUBLISH_ACK:
		c.handlePublishAck(pack)
	case SUBSCRIBE_ACK:
		c.handleSubscribeAck(pack)
	case UNSUBSCRIBE_ACK:
		c.handleUnSubscribeAck(pack)
	case DISCONNECT:
		c.handleDisconnect(pack)
	default:
		//忽略消息
	}
}

func (c *Client) handleConnect(pack *Pack) {

}

func (c *Client) handleConnectAck(pack *Pack) {

}

func (c *Client) handlePing(pack *Pack) {
	pack.Type = PONG
	err := c.Send(pack)
	if err != nil {

	}
}

func (c *Client) handlePong(pack *Pack) {

}

func (c *Client) handleRequest(pack *Pack) {
	if c.handler == nil {
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
	c.handler.ServeHTTP(rw, req)

	//3 回传 response
	pack.Type = RESPONSE
	pack.Payload = rw.Bytes()

	err = c.Send(pack)
	if err != nil {

	}
}

func (c *Client) handleResponse(pack *Pack) {

}

func (c *Client) handleStream(pack *Pack) {

}

func (c *Client) handleStreamEnd(pack *Pack) {

}

func (c *Client) handlePublish(pack *Pack) {

}

func (c *Client) handlePublishAck(pack *Pack) {

}

func (c *Client) handleSubscribeAck(pack *Pack) {

}

func (c *Client) handleUnSubscribeAck(pack *Pack) {

}

func (c *Client) handleDisconnect(pack *Pack) {

}
