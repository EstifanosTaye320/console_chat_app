package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func retriveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if (err!=nil) {
			fmt.Println("Error reading a message from the server")
			return
		}

		fmt.Println(message)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if (err!=nil) {
		fmt.Println("Error during the tcp handshake")
		return
	}
	defer conn.Close()

	go retriveMessage(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Conversation started type below: ")
	for {
		message, err := reader.ReadString('\n')
		if (err!=nil) {
			fmt.Println("Error reading the message")
			continue
		}

		message = strings.TrimSpace(message)

		if len(message) > 0 {
			fmt.Fprintln(conn, message)
		}
	}
}