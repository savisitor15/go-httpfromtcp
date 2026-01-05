package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/savisitor15/go-httpfromtcp/internal/request"
	"github.com/savisitor15/go-httpfromtcp/internal/response"
)

type Server struct {
	Running  *atomic.Bool
	listener *net.Listener
	handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handle Handler) (*Server, error) {
	if port > 65535 {
		return nil, fmt.Errorf("error: invalid port")
	}
	sPort := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", sPort)
	if err != nil {
		return nil, err
	}
	running := atomic.Bool{}
	running.Store(true)
	// create the server object
	s := Server{Running: &running, listener: &listener, handler: handle}
	// start the loop
	go s.listen()

	return &s, nil
}

func (s Server) Close() error {
	s.Running.Store(false)
	// get the listner
	cb := *s.listener
	err := cb.Close()
	return err
}

func (s Server) listen() {
	for {
		if s.Running.Load() == false {
			// server is dead
			break
		}
		listner := *s.listener
		conn, err := listner.Accept()
		if err != nil {
			if s.Running.Load() {
				fmt.Printf("error: openning incoming connection: %v\n", err)
			} else {
				break
			}
		}
		fmt.Println("handling connection")
		go s.handle(conn)
	}
}

func (s Server) handle(conn net.Conn) {
	defer conn.Close()
	if !s.Running.Load() {
		return
	}
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}
