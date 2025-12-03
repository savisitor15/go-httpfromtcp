package main

import (
	"fmt"
	"io"
	"os"
)

const INFILE = "messages.txt"

func main() {
	fp, err := os.Open(INFILE)
	if err != nil {
		fmt.Printf("error: opening file %s, %v", INFILE, err)
	}
	var buf []byte = make([]byte, 8)
	for {
		_, err := fp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("error: reading file: %v", err)
				break
			}
		}
		fmt.Printf("read: %s\n", buf)
	}
}
