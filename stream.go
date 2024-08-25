package pico

import "io"

func newStream(pico *Pico, id uint16) *stream {
	return &stream{
		c:    make(chan *Pack),
		pico: pico,
		id:   id,
	}
}

type stream struct {
	c    chan *Pack
	pico *Pico
	id   uint16

	buf []byte
}

func (s *stream) put(pack *Pack) {
	//todo 没有读操作怎么办
	s.c <- pack
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
	if len(s.buf) == 0 {
		pack := <-s.c
		if pack == nil {
			return 0, io.EOF
		}

		//如果是结束包，则
		if pack.Type == STREAM_END {
			close(s.c)
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

	//关闭管道
	close(s.c)

	//通知对方
	return s.pico.Send(&Pack{
		Id:      s.id,
		Type:    STREAM_END,
		Payload: nil,
	})
}
