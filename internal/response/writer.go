package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type writerState int

const (
	statInit writerState = iota
	stateStatusWritten
	stateHeadersWritten
	stateBodyWritten
)

type Writer struct {
	conn io.Writer
	state writerState
	Header headers.Headers
	status StatusCode
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		conn: conn,
		state: statInit,
		Header: headers.NewHeaders(),
	}
}

func (w *Writer) WriteStatusLine(code StatusCode) error {
	if w.state != statInit {
		return fmt.Errorf("status line already written or out of order")
	}
	w.status = code
	err := writeStatusLineTo(w.conn, code)
	if err == nil {
		w.state = stateStatusWritten
	}
	return err
}

func (w *Writer) WriteHeaders() error {
	if w.state != stateStatusWritten {
		return fmt.Errorf("must write status line before headers")
	}
	err := writeHeadersTo(w.conn, w.Header)
	if err == nil {
		w.state = stateHeadersWritten
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateHeadersWritten {
		return 0, fmt.Errorf("must write headers before body")
	}
	w.state = stateBodyWritten
	return w.conn.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateHeadersWritten {
		return 0, fmt.Errorf("must write headers before body")
	}
	_, err := fmt.Fprintf(w.conn, "%x\r\n", len(p))
	if err != nil {
		return 0, err
	}
	n, err := w.conn.Write(p)
	if err != nil {
		return n, err
	}
	_, err = w.conn.Write([]byte("\r\n"))
	if err != nil {
		return n, err
	}
	return n, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	err := w.WriteChunkedBodyDoneWithTrailers(headers.NewHeaders())
	return 0, err
}

func (w *Writer) WriteChunkedBodyDoneWithTrailers(h headers.Headers) error {
	if w.state != stateHeadersWritten {
		return fmt.Errorf("must write headers before finishing chunked body")
	}
	_, err := w.conn.Write([]byte("0\r\n"))
	if err != nil {
		return err
	}

	for k, v := range h {
		_, err := fmt.Fprintf(w.conn, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}

	_, err = w.conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.state = stateBodyWritten
	return nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.state != stateBodyWritten {
		return fmt.Errorf("must write body before trailers")
	}
	for k, v := range h {
		_, err := fmt.Fprintf(w.conn, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := w.conn.Write([]byte("\r\n"))
	return err
}