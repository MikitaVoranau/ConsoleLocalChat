package client

import (
	"CN_lab2/internal/additonal"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func StartClient() error {
	ip, port, err := additonal.GetIPandPorts()
	if err != nil {
		return fmt.Errorf("error getting address: %v", err)
	}
	var nickname string
	fmt.Println("Enter your nickname: ")
	_, err = fmt.Scanln(&nickname)
	if err != nil {
		return fmt.Errorf("error getting nickname: %v", err)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn.Close()
	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		text := consoleScanner.Text()
		if _, err := conn.Write([]byte(nickname + " > " + text + "\n")); err != nil {
			log.Printf("startclient : error %v", err)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading %v", err.Error())

		}
		fmt.Print("Server says: " + response)
		fmt.Println("Enter more text to send:")
	}

	if err := consoleScanner.Err(); err != nil {
		return fmt.Errorf("error reading from console %v", err.Error())
	}
	return nil
}
