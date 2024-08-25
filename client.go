package pico

import (
	"fmt"
)

type Client struct {
	Pico
}

func (c *Client) Connect() error {

	//login TODO 配置化
	auth := &Connect{Username: "", Password: ""}

	//auth TODO 支持Auth方式

	p := &Pack{Type: CONNECT, Content: auth}

	pack, err := c.Ask(p)
	if err != nil {
		return err
	}

	var ack ConnectAck
	err = pack.Decode(&ack)
	if err != nil {
		return err
	}

	if ack.Result {
		//todo 保存auth
	} else {
		//todo auth模式，改为login
	}

	return nil
}

func (c *Client) Disconnect(reason string) {
	_ = c.Send(&Pack{Type: DISCONNECT, Content: reason})
	_ = c.conn.Close()
}

func (c *Client) Publish(topic string, message any) error {
	pub := &Publish{Topic: topic, Message: message}
	pack, err := c.Ask(&Pack{Type: PUBLISH, Content: pub})
	if err != nil {
		return err
	}
	var ack PublishAck
	err = pack.Decode(&ack)
	if err != nil {
		return err
	}

	//逐一解析
	for t, r := range ack.Topics {
		if !r {
			return fmt.Errorf("topic %s not published", t)
		}
	}
	return nil
}

func (c *Client) Subscribe(filters []string) error {
	sub := &Subscribe{Filters: filters}
	pack, err := c.Ask(&Pack{Type: SUBSCRIBE, Content: sub})
	if err != nil {
		return err
	}
	var ack SubscribeAck
	err = pack.Decode(&ack)
	if err != nil {
		return err
	}

	//逐一解析
	for t, r := range ack.Filters {
		if !r {
			return fmt.Errorf("filter %s not subscribed", t)
		}
	}
	return nil
}

func (c *Client) Unsubscribe(filters []string) error {
	sub := &Unsubscribe{Filters: filters}
	pack, err := c.Ask(&Pack{Type: SUBSCRIBE, Content: sub})
	if err != nil {
		return err
	}
	var ack UnsubscribeAck
	err = pack.Decode(&ack)
	if err != nil {
		return err
	}

	//逐一解析
	for t, r := range ack.Filters {
		if !r {
			return fmt.Errorf("filter %s not unsubscribed", t)
		}
	}
	return nil
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
	//todo onMessage 回调
}
