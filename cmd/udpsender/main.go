package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const ADDRESS = "localhost:42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", ADDRESS)
	if err != nil {
		fmt.Printf("error: resolving %s using UDP: %v", ADDRESS, err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("error: opening up connection to %s: %v", ADDRESS, err)
		return
	}
	defer conn.Close()
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		str, err := buf.ReadString('\n')
		if err != nil {
			fmt.Printf("error: reading input: %v\n", err)
			continue
		}
		written, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Printf("error: writing to the connection socket: %v\n", err)
			continue
		}
		fmt.Printf("\nWrote %d bytes to %s\n", written, ADDRESS)
	}
}
