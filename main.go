package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
			fmt.Printf("error: openning incomminc connection: %v", err)
		}
		producer := getLinesChannel(conn)
		for val := range producer {
			fmt.Println(val)
		}
	}
}

func getLinesChannel(fp io.ReadCloser) <-chan string {
	var out string
	ch := make(chan string)
	go func(fp io.ReadCloser) {
		for {
			buf := make([]byte, 8)
			n, err := fp.Read(buf)
			if err != nil {
				if err == io.EOF {
					if len(out) != 0 {
						ch <- out
					}
					close(ch)
					//fp.Close()
					break
				} else {
					fmt.Printf("error: reading file: %v", err)
					break
				}
			}
			parts := strings.Split(string(buf[0:n]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s", out, parts[i])
				out = ""
			}
			out += parts[len(parts)-1]
		}
	}(fp)
	return ch
}
