package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state requestState
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const crlf = "\r\n"

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	r := &Request{state: requestStateInitialized}

	for r.state != requestStateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf) * 2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		readToIndex += n

		parsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if parsed > 0 {
			copy(buf, buf[parsed:readToIndex])
			readToIndex -= parsed
		}
	}

	if r.state != requestStateDone {
		return nil, fmt.Errorf("incomplete request: parser did not reach done state")
	}

	return r, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + len(crlf), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request line: expected 3 parts, got %d", len(parts))
	}

	method := parts[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2{
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecgnized HTTP-version: %s", version)
	}

	return &RequestLine{
		HttpVersion: versionParts[1],
		RequestTarget: requestTarget,
		Method: method,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.state = requestStateDone
		return n, nil

	case requestStateDone:
		return 0, fmt.Errorf("error: tyring to read data in a done state")

	default:
		return 0, fmt.Errorf("error: unknown state %d", r.state)
	}
}