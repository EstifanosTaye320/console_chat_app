package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func retriveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if (err!=nil) {
			fmt.Println("Error reading a message from the server")
			continue
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
	for {
		fmt.Print("Write a message for the group: ")
		message, err := reader.ReadString('\n')
		if (err!=nil) {
			fmt.Println("Error reading the message")
			continue
		}

		fmt.Fprintln(conn, message)
	}
}