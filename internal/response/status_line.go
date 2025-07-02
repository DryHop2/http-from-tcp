package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusOK 					StatusCode = 200
	StatusBadRequest 			StatusCode = 400
	StatusInternalServerError 	StatusCode = 500
)

func getStatusLine(statusCode StatusCode) string {
	var reason string

	switch statusCode {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "Internal Server Error"
	default:
		reason = ""
	}
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason)
}

func writeStatusLineTo(w io.Writer, statusCode StatusCode) error {
	_, err := fmt.Fprint(w, getStatusLine(statusCode))
	return err
}
