package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			continue
		}
		fmt.Println("Connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
	fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <- chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		var currentLine string
		buffer := make([]byte, 8)
		
		for {
			bytesRead, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				break
			}
			chunk := string(buffer[:bytesRead])
			parts := strings.Split(chunk, "\n")
			for i, part := range parts {
				if i < len(parts) - 1 {
					lines <- currentLine + part
					currentLine = ""
				} else {
					currentLine += part
				}
			}
		}
	}()
	return lines
}