package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"io"
	"log"
	"net"
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

		go handleConnection(conn)
	}
}

func handleConnection(conn io.ReadCloser) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Failed to parse request:", err)
		return
	}

	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- Version:", req.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, value := range req.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	}

	fmt.Println("Connection closed")
}