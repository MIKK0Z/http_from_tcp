package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const PORT = ":42069"

func main() {

	file, _ := os.Open("main.go")
	f_channel := getLinesChannel(file)
	for line := range f_channel {
		fmt.Println(line)
	}

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Listening for TCP on %s\n", PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		channel := getLinesChannel(conn)
		for msg := range channel {
			fmt.Printf("read: %s\n", msg)
		}

		fmt.Printf("Close connection to %s\n", conn.RemoteAddr())
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)
	current_line := ""

	go func() {
		defer close(channel)
		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)

			if err != nil {
				if current_line != "" {
					channel <- current_line
					current_line = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			parts := strings.Split(string(buff[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				current_line += parts[i]
				channel <- current_line
				current_line = ""
			}

			current_line += parts[len(parts)-1]
		}
	}()

	return channel
}
