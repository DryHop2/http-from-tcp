package server

import (
	"fmt"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed atomic.Bool
}

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{listener: ln}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	const body = ""
	status := response.StatusOK
	headers := response.GetDefaultHeaders(len(body))

	if err := response.WriteStatusLine(conn, status); err != nil {
		log.Printf("error writing status line: %v", err)
		return
	}
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Printf("error writing headers: %v", err)
		return
	}
	if _, err := fmt.Fprint(conn, "\r\n"); err != nil {
		log.Printf("error writing CRLF: %v", err)
		return
	}
}