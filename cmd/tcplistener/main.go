package main

import (
	"fmt"
	"net"

	"github.com/savisitor15/go-httpfromtcp/internal/request"
)

const INFILE = "messages.txt"
const NETWORK = "tcp"
const ADDRESS = ":42069"

func main() {
	listener, err := net.Listen(NETWORK, ADDRESS)
	if err != nil {
		fmt.Printf("error: opening port %s, %v", ADDRESS, err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error: openning incoming connection: %v", err)
		}
		producer, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error: opening reader: %v", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", producer.RequestLine.Method)
		fmt.Printf("- Target: %s\n", producer.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", producer.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for token, value := range producer.Headers {
			fmt.Printf("- %s: %s\n", token, value)
		}

	}
}
