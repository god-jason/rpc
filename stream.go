package pico

import "io"

func newStream(pico *Pico, id uint16) *stream {
	return &stream{
		reader: make(chan *Pack),
		pico:   pico,
		id:     id,
	}
}

type stream struct {
	reader chan *Pack
	pico   *Pico
	id     uint16

	buf []byte
}

func (s *stream) put(pack *Pack) {
	s.reader <- pack
}

func (s *stream) Write(buf []byte) (int, error) {
	pack := &Pack{
		Id:      s.id,
		Type:    STREAM,
		Payload: buf,
	}
	return len(buf), s.pico.Send(pack)
}

func (s *stream) Read(buf []byte) (int, error) {
	//阻塞读数据
	if len(buf) == 0 {
		pack := <-s.reader
		if pack == nil {
			return 0, io.EOF
		}

		s.buf = pack.Payload //复制
		//s.buf = make([]byte, len(pack.Payload))
		//copy(s.buf, pack.Payload)
	}

	n := copy(buf, s.buf)
	if n == len(s.buf) {
		return n, nil
	}

	//保存剩余
	s.buf = s.buf[n:]
	return n, nil
}

func (s *stream) Close() error {
	close(s.reader)
	return nil
}
