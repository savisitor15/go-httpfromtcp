package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/savisitor15/go-httpfromtcp/internal/request"
	"github.com/savisitor15/go-httpfromtcp/internal/response"
	"github.com/savisitor15/go-httpfromtcp/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	herr := server.HandlerError{}
	const interrpath = "/myproblem"
	const clienterrpath = "/yourproblem"
	const goodresp = "All good, frfr\n"
	const cerrresp = "Your problem is not my problem\n"
	const ierrresp = "Woopsie, my bad\n"
	target := req.RequestLine.RequestTarget

	if strings.HasSuffix(target, clienterrpath) {
		herr.Code = response.StatusCodeBadRequest
		herr.Message = cerrresp
		return &herr
	}
	if strings.HasSuffix(target, interrpath) {
		herr.Code = response.StatusCodeInternal
		herr.Message = ierrresp
		return &herr
	}
	// all good?
	herr.Code = response.StatusCodeOK
	w.Write([]byte(goodresp))
	return &herr
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
