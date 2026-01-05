package response

import "fmt"

type StatusCode int

const ServerVersion = "HTTP/1.1"

const (
	StatusCodeOK         StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternal   StatusCode = 500
)

func ResponseStatusString(code int) string {
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

func getStatusLine(statusCode StatusCode) []byte {
	return []byte(fmt.Sprintf("%s %d %s\r\n", ServerVersion, statusCode, ResponseStatusString(int(statusCode))))
}
