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

	line := data[:idx]
	parts := bytes.SplitN(line, []byte(":"), 2)
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("malformed header line (no colon): %q", line)
	}

	rawKey := string(parts[0])
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
	value := strings.TrimSpace(string(parts[1]))
	h.Set(key, value)

	return idx + len(crlf), false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	if existing, ok := h[key]; ok {
		h[key] = existing + ", " + value
	} else {
		h[key] = value
	}
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	return h[key]
}

func isValidHeaderField(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 128 || !tcharTable[s[i]] {
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	delete(h, strings.ToLower(key))
}