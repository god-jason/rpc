package pico

import (
	"errors"
)

type Pack struct {
	Id       uint16
	Type     uint8
	Encoding uint8
	Payload  []byte
	Content  any
}

func (p *Pack) Encode() (err error) {
	switch v := p.Content.(type) {
	case nil:
	case string:
		p.Payload = []byte(v)
	case []byte:
		p.Payload = v
	default:
		if p.Encoding == BINARY {
			p.Encoding = JSON //默认使用JSON
		}
		if encoder, ok := encoders[p.Encoding]; ok {
			p.Payload, err = encoder(v)
		} else {
			err = errors.New("encoding not supported")
		}
	}
	return
}

func (p *Pack) Decode(payload any) (err error) {
	if p.Encoding == BINARY {
		return nil
	}

	if decoder, ok := decoders[p.Encoding]; ok {
		err = decoder(p.Payload, payload)
	} else {
		err = errors.New("encoding not supported")
	}
	return
}
