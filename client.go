package pico

type Client struct {
	Pico
}

func (c *Client) Connect() error {
	auth := &Connect{
		Username: "",
		Password: "",
		Id:       "",
		Token:    "",
	}
	c.Ask(&Pack{Type: CONNECT, Encoding: JSON, Payload: nil})
}

func (c *Client) Disconnect(reason string) {
	_ = c.Send(&Pack{Type: DISCONNECT, Payload: []byte(reason)})
	_ = c.conn.Close()
}

func (c *Client) Publish(pack *Pack) error {

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

func (c *Client) handle(pack *Pack) {
	switch pack.Type {
	case DISCONNECT:
		c.handleDisconnect(pack)
	case CONNECT_ACK:
		c.handleAsk(pack)
	case PING:
		c.handlePing(pack)
	case PONG:
		c.handleAsk(pack)
	case REQUEST:
		c.handleRequest(pack)
	case RESPONSE:
		c.handleAsk(pack)
	case STREAM, STREAM_END:
		c.handleStream(pack)
	case PUBLISH_ACK, SUBSCRIBE_ACK, UNSUBSCRIBE_ACK:
		c.handleAsk(pack)
	case MESSAGE:
		c.handleMessage(pack)
	default:
		//忽略消息
	}
}

func (c *Client) handleDisconnect(pack *Pack) {
	_ = c.conn.Close()
}

func (c *Client) handleMessage(pack *Pack) {
	//onMessage
}
