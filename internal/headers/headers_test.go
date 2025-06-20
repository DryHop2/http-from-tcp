package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func NewHeaders() Headers {
	return make(Headers)
}

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["User-Agent"] = "curl"
	data = []byte("Accept: */*\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "curl", headers["User-Agent"])
	assert.Equal(t, "*/*", headers["accept"])
	assert.Equal(t, len(data), n)
	assert.False(t, done)

	// Test: Valid done (CRLF only)
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header (space before colon)
	headers = NewHeaders()
	data = []byte("  Host : localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Capital letters in header keys
	headers = NewHeaders()
	data = []byte("HOST: LOCALhoST:42069\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "LOCALhoST:42069", headers["host"])
	assert.Equal(t, len(data), n)
	assert.False(t, done)

	// Test: Multiple values for single header key
	headers = NewHeaders()
	headers["host"] = "localhost:42069"
	data = []byte("Host: localhost:33333\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069, localhost:33333", headers["host"])
	assert.Equal(t, len(data), n)
	assert.False(t, done)
}