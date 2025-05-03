package client

import (
	"CN_lab2/internal/additonal"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func StartClient() error {
	ip, port, err := additonal.GetIPandPorts()
	if err != nil {
		return fmt.Errorf("error getting address: %v", err)
	}

	var nickname string
	fmt.Println("Enter your nickname: ")
	_, err = fmt.Scanln(&nickname)
	fmt.Println()
	if err != nil {
		return fmt.Errorf("error getting nickname: %v", err)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return fmt.Errorf("error connection: %w", err)
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				fmt.Println(line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error scanning: %v", err)
		}
	}()

	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		text := consoleScanner.Text()
		if _, err := conn.Write([]byte(nickname + " (" + time.Now().Format("15:04:05") + ") -> " + text + "\n")); err != nil {
			log.Printf("startclient : error %v", err)
		}
	}
	if err := consoleScanner.Err(); err != nil {
		return fmt.Errorf("error reading from console %v", err.Error())
	}
	return nil
}
