package pico

import (
	"encoding/xml"
	"github.com/bytedance/sonic"
	"github.com/shamaton/msgpack/v2"
	"gopkg.in/yaml.v3"
)

type Encoder func(any) ([]byte, error)
type Decoder func([]byte, any) error

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
