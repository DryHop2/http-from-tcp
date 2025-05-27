package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("read: %s\n", buffer[:bytesRead])
	}

}