package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func (h Headers) Get(token string) (string, bool) {
	// case in-sensitive
	token = strings.ToLower(token)
	value, found := h[token]
	return value, found
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	outLen := idx + 2
	if strings.HasPrefix(string(data), crlf) {
		// blank line
		return 2, true, nil
	}
	newHead, err := headerFromString(string(data[:idx]))
	if err != nil {
		return 0, false, err
	}
	h.Set(newHead[0], newHead[1])
	return outLen, false, nil
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	delete(h, key)
}

func headerFromString(s string) ([]string, error) {
	// Left trim first
	s = strings.TrimLeft(s, " ")
	re := regexp.MustCompile(`^(?P<token>[0-9A-Za-z\!\#\$\%\&\'\*\+-\.\^\_\x60\|\~]+)(:)(?P<space>\s+)(?P<value>.+)`)
	result := make(map[string]string)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil, fmt.Errorf("invalid header token")
	}
	for i, name := range re.SubexpNames() {
		result[name] = match[i]
	}
	return []string{strings.ToLower(result["token"]), result["value"]}, nil
}

func NewHeaders() Headers {
	return Headers{}
}
