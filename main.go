package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const INFILE = "messages.txt"

func main() {
	fp, err := os.Open(INFILE)
	if err != nil {
		fmt.Printf("error: opening file %s, %v", INFILE, err)
	}
	producer := getLinesChannel(fp)
	for val := range producer {
		fmt.Printf("read: %s\n", val)
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
					fp.Close()
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
