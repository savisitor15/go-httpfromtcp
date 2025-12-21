package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"
const sep = ":"
const space = " "

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	outLen := idx + 2
	if strings.HasPrefix(string(data), crlf) {
		// blank line
		return outLen, true, nil
	}
	newHead, err := headerFromString(string(data[:idx]))
	if err != nil {
		return 0, false, err
	}
	h[newHead[0]] = newHead[1]
	return outLen, false, nil
}

func headerFromString(s string) ([]string, error) {
	// split the header
	parts := strings.SplitN(s, sep, 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid header line")
	}
	// get the key and ensure validity - space may be left trimmed but not right
	head := strings.TrimLeft(parts[0], space)
	if strings.HasSuffix(head, " ") {
		return nil, fmt.Errorf("invalid header key")
	}
	// get value, trim leading spaces
	value := strings.TrimLeft(parts[1], space)
	return []string{head, value}, nil
}

func NewHeaders() Headers {
	return Headers{}
}
