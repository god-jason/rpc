package pico

import (
	"fmt"
	"net"
)

type Server struct {
}

func (s *Server) Serve(port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		in := &Incoming{Pico: Pico{conn: conn}, server: s}
		in.init()

		go in.receive()
	}

}
