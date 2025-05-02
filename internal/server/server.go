package server

import (
	"CN_lab2/internal/additonal"
	"bufio"
	"fmt"
	"log"
	"net"
)

func StartServer() error {
	ip, port, err := additonal.GetIPandPorts()
	if err != nil {
		return fmt.Errorf("error getting address: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server started on %s:%d\n", ip, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		clientMessage := scanner.Text()
		fmt.Printf("%s\n", clientMessage)

		if _, err := conn.Write([]byte("Message received.\n")); err != nil {
			log.Printf("Error writing to client: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("the user disconnected: %v", err)
	}
}
