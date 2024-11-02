package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var lstClient []net.Conn
var mu sync.Mutex

func broadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	for _, conn := range lstClient {
		fmt.Fprintln(conn, message)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if (err!=nil) {
			fmt.Println("Error reading a message", err)
			continue
		}

		go broadcastMessage(message)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if (err!=nil) {
		fmt.Println("Error createing the sockate connection", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if (err!=nil) {
			fmt.Println("Error during the tcp handshake", err)
			continue
		}

		mu.Lock()
		lstClient = append(lstClient, conn)
		mu.Unlock()

		go handleRequest(conn)
	}
}