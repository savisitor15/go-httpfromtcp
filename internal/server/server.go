package server

import (
	"bytes"
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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		// 400
		fmt.Printf("error handling request: %v", err)
		err = response.WriteStatusLine(conn, response.StatusCodeBadRequest)
		if err != nil {
			fmt.Printf("error to report error!")
		}
		return
	}
	buf := new(bytes.Buffer)
	hErr := s.handler(buf, req)
	if hErr.Code != response.StatusCodeOK {
		header := response.GetDefaultHeaders(len(hErr.Message))
		err = response.WriteStatusLine(conn, hErr.Code)
		if err != nil {
			fmt.Printf("error to report error status line!")
		}
		err = response.WriteHeaders(conn, header)
		if err != nil {
			fmt.Printf("error to report error headers!")
		}
		_, err = conn.Write([]byte(hErr.Message))
		if err != nil {
			fmt.Printf("error to report error message!")
		}
	}
	header := response.GetDefaultHeaders(len(buf.Bytes()))
	err = response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		fmt.Printf("error writing status line: %v", err)
	}
	err = response.WriteHeaders(conn, header)
	if err != nil {
		fmt.Printf("error writing headers: %v", err)
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("error writing body: %v", err)
	}
}
