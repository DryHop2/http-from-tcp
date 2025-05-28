package main

import (
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

	buffer := make([]byte, 8)
	var currentLine string

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		chunk := string(buffer[:bytesRead])
		parts := strings.Split(chunk, "\n")
		for i, part := range parts {
			if i < len(parts)-1 {
				fmt.Printf("read: %s%s\n", currentLine, part)
				currentLine = ""
			} else {
				currentLine += part
			}
		}
	}
	fmt.Printf("read: %s", currentLine)

}