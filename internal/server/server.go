package server

import (
	"CN_lab2/internal/additonal"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn     net.Conn
	nickname string
}

var (
	clients = make(map[net.Conn]Client)
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		log.Printf("Failed to read nickname from client")
		return
	}
	nickname := scanner.Text()
	client := Client{conn: conn, nickname: nickname}

	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	log.Printf("%s connected", nickname)
	broadcastMessage(fmt.Sprintf("%s connected to chat\n", nickname), conn)

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		log.Printf("%s disconnected", nickname)
		broadcastMessage(fmt.Sprintf("%s disconnected from chat\n", nickname), conn)
	}()

	for scanner.Scan() {
		msg := scanner.Text()

		if strings.HasPrefix(msg, "FILE:") {
			handleFileMessage(conn, msg)
		} else {
			broadcastMessage(fmt.Sprintf("%s (%s) -> %s\n", client.nickname, time.Now().Format("15:04:05"), msg), conn)
		}
	}
}
func broadcastMessage(msg string, excludeConn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range clients {
		if client.conn != excludeConn {
			_, err := client.conn.Write([]byte(msg))
			if err != nil {
				log.Printf("Error sending message to client %s: %v", client.nickname, err)
				client.conn.Close()
				delete(clients, client.conn)
			}
		}
	}
}

func handleFileMessage(conn net.Conn, msg string) {
	msg = strings.TrimPrefix(msg, "FILE:")
	parts := strings.SplitN(msg, "|", 4)
	if len(parts) < 4 {
		log.Printf("Invalid file metadata format")
		return
	}

	sender := parts[0]
	timestamp := parts[1]
	filename := parts[2]
	filesize, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil || filesize <= 0 {
		log.Printf("Invalid file size: %v", err)
		return
	}

	// Читаем ровно filesize байт
	fileData := make([]byte, filesize)
	_, err = io.ReadFull(conn, fileData)
	if err != nil {
		log.Printf("Error reading file data: %v", err)
		return
	}

	broadcastFile(sender, timestamp, filename, filesize, fileData)
}

func broadcastFile(sender, timestamp, filename string, filesize int64, fileData []byte) {
	meta := fmt.Sprintf("FILE:%s|%s|%s|%d\n", sender, timestamp, filename, filesize)

	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range clients {
		_, err := client.conn.Write([]byte(meta))
		if err != nil {
			log.Printf("Ошибка отправки метаданных клиенту %s: %v", client.nickname, err)
			continue
		}
		_, err = client.conn.Write(fileData)
		if err != nil {
			log.Printf("Ошибка отправки данных клиенту %s: %v", client.nickname, err)
			client.conn.Close()
			delete(clients, client.conn)
		}
	}
}
