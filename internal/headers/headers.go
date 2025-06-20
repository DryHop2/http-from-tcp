package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

var tcharTable [256]bool

func init() {
	for c := 'a'; c <= 'z'; c++ {
		tcharTable[c] = true
	}
	for c := 'A'; c <= 'Z'; c++ {
		tcharTable[c] = true
	}
	for c := '0'; c <= '9'; c++ {
		tcharTable[c] = true
	}
	specials := "!#$%&'*+-.^_`|~"
	for i := 0; i < len(specials); i++ {
		tcharTable[specials[i]] = true
	}
}

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

	rawKey := line[:colonIdx]
	if strings.TrimSpace(rawKey) != rawKey {
		return 0, false, fmt.Errorf("invalid header: space before colon")
	}
	if rawKey == "" {
		return 0, false, fmt.Errorf("empty header key")
	}
	if !isValidHeaderField(rawKey) {
		return 0, false, fmt.Errorf("invalid header: bad character in %q", rawKey)
	}

	key := strings.ToLower(rawKey)
	value := strings.TrimSpace(line[colonIdx + 1:])
	currentValue, exists := h[key]
	if exists {
		h[key] = currentValue + ", " + value
	} else {
		h[key] = value
	}

	return idx + len(crlf), false, nil
}

func isValidHeaderField(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 128 || !tcharTable[s[i]] {
			return false
		}
	}
	return true
}