package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed atomic.Bool
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: ln,
		handler: handler,
	}
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		writer := response.NewWriter(conn)
		writer.WriteStatusLine(response.StatusBadRequest)
		writer.Header.Set("Content-Type", "text/html")
		writer.WriteHeaders()
		writer.WriteBody([]byte(`<html>
									<head>
										<title>400 Bad Request</title>
									</head>
									<body>
										<h1>Bad Request</h1>
										<p>Your request honestly kinda sucked.</p>
									</body>
									</html>`))
		return
	}

	writer := response.NewWriter(conn)
	s.handler(writer, req)
}