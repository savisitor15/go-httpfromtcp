package response

import (
	"fmt"

	"github.com/savisitor15/go-httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	header.Set("Connection", "Close")
	header.Set("Content-Type", "text/plain")
	return header
}
