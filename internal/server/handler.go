package server

import (
	"io"

	"github.com/savisitor15/go-httpfromtcp/internal/request"
	"github.com/savisitor15/go-httpfromtcp/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
