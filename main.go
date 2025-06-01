package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
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