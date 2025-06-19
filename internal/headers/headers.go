package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return len(crlf), true, nil
	}

	line := string(data[:idx])
	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("malformed header line (no colon): %q", line)
	}

	parts := strings.SplitN(string(line), ":", 2)
	if strings.TrimSpace(parts[0]) != parts[0] {
		return 0, false, fmt.Errorf("invalid header: space before colon")
	}
	
	key := strings.TrimSpace(line[:colonIdx])
	value := strings.TrimSpace(line[colonIdx + 1:])

	if key == "" {
		return 0, false, fmt.Errorf("empty header key in line: %q", line)
	}

	h[key] = value

	return idx + len(crlf), false, nil
}

