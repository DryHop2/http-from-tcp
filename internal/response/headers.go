package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)

	strContentLen := strconv.Itoa(contentLen)
	h.Set("Content-Length", strContentLen)
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func writeHeadersTo(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\r\n")
	return err
}