package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const ADDR = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ADDR)
	if err != nil {
		os.Exit(1)
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		os.Exit(1)
	}
	defer udpConn.Close()

	fmt.Printf("Sending to %s\n", ADDR)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}

		_, err = udpConn.Write([]byte(msg))
		if err != nil {
			os.Exit(1)
		}
	}
}
