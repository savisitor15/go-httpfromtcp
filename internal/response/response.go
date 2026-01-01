package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/savisitor15/go-httpfromtcp/internal/headers"
)

type StatusCode int

const ServerVersion = "HTTP/1.1"

const (
	StatusCodeOK         StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternal   StatusCode = 500
)

func StatusString(code int) string {
	switch code {
	case int(StatusCodeOK):
		return "OK"
	case int(StatusCodeBadRequest):
		return "Bad Request"
	case int(StatusCodeInternal):
		return "Internal Server Error"
	default:
		return "Unkown Status"
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	conLenStr := strconv.Itoa(contentLen)
	header.Set("Content-Length", conLenStr)
	header.Set("Connection", "Close")
	header.Set("Content-Type", "text/plain")
	return header
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(
		[]byte(
			fmt.Sprintf("%s %d %s\r\n", ServerVersion, statusCode, StatusString(int(statusCode)))))
	return err
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var err error
	for token, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", token, value)))
		if err != nil {
			return err
		}
	}
	// close down the headers
	_, err = w.Write([]byte("\r\n"))
	return err
}
