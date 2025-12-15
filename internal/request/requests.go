package request

import (
	"fmt"
	"io"
	"regexp"

	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func validateVerbs(method string) (string, error) {
	validVerbs := [8]string{"OPTIONS",
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"TRACE",
		"CONNECT"}
	err := fmt.Errorf("error: invalid method")
	if len(method) < 3 {
		// shorter than the shortest verbs
		return "", err
	}
	for i := range validVerbs {
		if validVerbs[i] == method {
			return method, nil
		}
	}
	return "", err
}

func validateVersion(raw string) (string, error) {
	err := fmt.Errorf("error: version format invalid")
	// RFC validation
	valid := regexp.MustCompile(`HTTP\/\d\.\d`).MatchString(raw)
	if !valid {
		return "", err
	}
	ver, found := strings.CutPrefix(raw, "HTTP/")
	if found != true {
		return "", err
	}
	return ver, nil

}

func parseRequestLine(raw string) (*RequestLine, error) {
	lines := strings.Split(raw, "\r\n")
	if len(lines) <= 0 {
		return nil, fmt.Errorf("error: invalid string format")
	}
	// first element should be the reuqest line
	rawRL := strings.Split(lines[0], " ")
	if len(rawRL) < 3 {
		return nil, fmt.Errorf("error: invalid request line format")
	}
	ver, err := validateVersion(rawRL[2])
	if err != nil {
		return nil, err
	}
	method, err := validateVerbs(rawRL[0])
	if err != nil {
		return nil, err
	}
	return &RequestLine{
		Method:        method,
		RequestTarget: rawRL[1],
		HttpVersion:   ver,
	}, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reqline, err := parseRequestLine(string(raw))
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *reqline}, nil
}
