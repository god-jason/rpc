package pico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/http"
)

type Client struct {
	handler http.Handler
	conn    net.Conn
	id      uint16
}

func (c *Client) Disconnect(reason string) {
	//todo 发送 disconnect
	_ = c.conn.Close()
}

func (c *Client) receive() {
	header := make([]byte, HeaderSize)
	buf := make([]byte, BufferSize)

	//bufio.NewReaderSize(c.conn, BufferSize)
	reader := bufio.NewReader(c.conn)

	for {
		n, e := reader.Read(header)
		if e != nil {
			break
		}
		if n < HeaderSize {
			//break
			continue //TODO继续接受
			//reader.UnreadByte()
		}

		if bytes.Compare(header[:MagicSize], []byte(Magic)) != 0 {
			continue
		}

		//解析
		pack := &Pack{
			Id:       binary.BigEndian.Uint16(header[4:]),
			Type:     header[6] >> 4,
			Encoding: header[6] & 0x0f,
		}

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
				//TODO disconnect
				break
			}

			//todo 内存复制问题
			pack.Payload = b[:n]
		}

		c.handle(pack)
	}
}

func (c *Client) Send(pack *Pack) error {
	header := make([]byte, HeaderSize)
	copy(header[:MagicSize], []byte(Magic))
	binary.BigEndian.PutUint16(header[4:], c.id)
	header[6] = pack.Type<<4 + pack.Encoding
	length := len(pack.Payload)
	header[7] = byte(length >> 16)
	header[8] = byte(length >> 8)
	header[9] = byte(length)

	//todo 一次写完

	_, err := c.conn.Write(header)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(pack.Payload)
	if err != nil {
		return err
	}

	return nil
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
