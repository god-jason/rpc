package pico

type Incoming struct {
	Pico
	server *Server
}

func (c *Incoming) receive() {
	for {
		pack, err := c.readPack()
		if err != nil {
			break
		}
		c.handle(pack)
	}
}

func (c *Incoming) handle(pack *Pack) {
	switch pack.Type {
	case DISCONNECT:
		c.handleDisconnect(pack)
	case CONNECT:
		c.handleConnect(pack)
	case PING:
		c.handlePing(pack)
	case PONG:
		c.handleAck(pack)
	case REQUEST:
		c.handleRequest(pack)
	case RESPONSE:
		c.handleAck(pack)
	case STREAM, STREAM_END:
		c.handleStream(pack)
	case PUBLISH:
		c.handlePublish(pack)
	case PUBLISH_ACK:
		c.handleAck(pack)
	case SUBSCRIBE:
		c.handleSubscribe(pack)
	case UNSUBSCRIBE:
		c.handleUnSubscribe(pack)
	default:
		//忽略消息
	}
}

func (c *Incoming) Disconnect(reason string) {
	_ = c.Send(&Pack{Type: DISCONNECT, Content: reason})
	_ = c.conn.Close()
}

func (c *Incoming) handleConnect(pack *Pack) {

	//todo 鉴权

}

func (c *Incoming) handlePublish(pack *Pack) {

}

func (c *Incoming) handleSubscribe(pack *Pack) {

}

func (c *Incoming) handleUnSubscribe(pack *Pack) {

}

func (c *Incoming) handleDisconnect(pack *Pack) {
	_ = c.conn.Close()
}
