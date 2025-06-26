package server

import (
	"bytes"
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
		herr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message: "Invalid Request\n",
		}
		WriteHandlerError(conn, herr)
		return
	}

	var respBody bytes.Buffer

	herr := s.handler(&respBody, req)
	if herr != nil {
		WriteHandlerError(conn, herr)
		return
	}

	// bodyBytes := respBody.Bytes()
	headers := response.GetDefaultHeaders(respBody.Len())

	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Printf("failed to write status line: %v", err)
		return
	}
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Printf("failed to write headers: %v", err)
		return
	}
	if _, err := respBody.WriteTo(conn); err != nil {
		log.Printf("failed to write body: %v", err)
	}
}