package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

var lstClient []net.Conn
var mu sync.Mutex

func loadMessages(conn net.Conn, path string) {
	dat, _ := os.ReadFile(path)
	lines := strings.Split(string(dat), "\n")
	for _, line := range lines {
		fmt.Fprintln(conn, line)
	}
}

func broadcastMessage(message string, conn net.Conn, f *os.File) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Broadcasting ", message)
	for _, hconn := range lstClient {
		if hconn!=conn {
			fmt.Fprintln(hconn, message)
		}
	}

	f.WriteString(message + "\n")
}

func handleRequest(conn net.Conn, f *os.File) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if (err!=nil) {
			fmt.Println("Error reading a message", err)
			return
		}

		message = strings.TrimSpace(message)

		if len(message) > 0 {
			broadcastMessage(message, conn, f)
		}
	}
}

func main() {
	path, _ := filepath.Abs("data/conversation.txt")
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	listener, err := net.Listen("tcp", ":8080")
	if (err!=nil) {
		fmt.Println("Error createing the sockate connection", err)
		return
	}

	defer func () {
		f.Write([]byte(""))
		f.Close()
		listener.Close()
	}() 
	fmt.Println("Server is running of port 8080...")
	
	go func () {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		
		f.Close()
		os.Remove(path)
		listener.Close()
		os.Exit(1)
	}()

	for {
		conn, err := listener.Accept()
		if (err!=nil) {
			fmt.Println("Error during the tcp handshake", err)
			continue
		}

		mu.Lock()
		lstClient = append(lstClient, conn)
		mu.Unlock()

		loadMessages(conn, path)

		go handleRequest(conn, f)
	}
}