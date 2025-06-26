package server

import (
	"fmt"
	"io"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode  response.StatusCode
	Message		string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, herr *HandlerError) {
	headers := response.GetDefaultHeaders(len(herr.Message))
	_ = response.WriteStatusLine(w, herr.StatusCode)
	_ = response.WriteHeaders(w, headers)
	fmt.Fprint(w, herr.Message)
}