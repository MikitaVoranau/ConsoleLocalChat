package server

import (
	"CN_lab2/internal/additonal"
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	clients = make(map[net.Conn]bool)
	mutex   = &sync.Mutex{}
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

		mutex.Lock()
		clients[conn] = true
		mutex.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
	}()

	for scanner.Scan() {
		clientMessage := scanner.Text()
		fmt.Printf("Received: %s\n", clientMessage)
		mutex.Lock()
		for client := range clients {
			if client != conn {
				_, err := client.Write([]byte(clientMessage + "\n"))
				if err != nil {
					log.Printf("Error sending message to client: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("the user disconnected: %v", err)
	}
}
