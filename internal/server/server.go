package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	Running  *atomic.Bool
	listener *net.Listener
}

func Serve(port int) (*Server, error) {
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
	s := Server{Running: &running, listener: &listener}
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
	const response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"
	if !s.Running.Load() {
		return
	}
	conn.Write([]byte(response))
	defer conn.Close()
}
