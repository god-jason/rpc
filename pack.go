package pico

import (
	"encoding/xml"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/shamaton/msgpack/v2"
	"gopkg.in/yaml.v3"
)

var ErrEncoding = errors.New("编码不支持")
var ErrNotEnough = errors.New("长度不足")

type Encoder func(any) ([]byte, error)
type Decoder func([]byte, any) error

const (
	DISCONNECT uint8 = iota
	CONNECT
	CONNECT_ACK
	PING
	PONG
	REQUEST
	RESPONSE
	STREAM
	STREAM_END
	PUBLISH
	PUBLISH_ACK
	SUBSCRIBE
	SUBSCRIBE_ACK
	UNSUBSCRIBE
	UNSUBSCRIBE_ACK
	MESSAGE
)

const (
	BINARY uint8 = iota
	JSON
	XML
	YAML
	MSGPACK
)

var encoders = map[uint8]Encoder{
	JSON:    sonic.Marshal,
	XML:     xml.Marshal,
	YAML:    yaml.Marshal,
	MSGPACK: msgpack.Marshal,
}

var decoders = map[uint8]Decoder{
	JSON:    sonic.Unmarshal,
	XML:     xml.Unmarshal,
	YAML:    yaml.Unmarshal,
	MSGPACK: msgpack.Unmarshal,
}

func RegisterEncoding(typ uint8, encoder Encoder, decoder Decoder) {
	encoders[typ] = encoder
	decoders[typ] = decoder
}

type Pack struct {
	Id       uint16
	Type     uint8
	Encoding uint8
	Payload  []byte
}
