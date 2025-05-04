package client

import (
	"CN_lab2/internal/additonal"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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
	if err != nil {
		return fmt.Errorf("error getting nickname: %v", err)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(nickname + "\n")); err != nil {
		return fmt.Errorf("error sending nickname: %v", err)
	}

	go func() {
		reader := bufio.NewReader(conn)
		for {
			header, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Connection closed: %v", err)
				return
			}

			if strings.HasPrefix(header, "FILE:") {
				fileMeta := strings.TrimPrefix(header, "FILE:")
				parts := strings.SplitN(fileMeta, "|", 4)
				if len(parts) < 4 {
					continue
				}

				sender := parts[0]
				timestamp := parts[1]
				filename := parts[2]
				filesize, _ := strconv.ParseInt(parts[3], 10, 64)

				fileData := make([]byte, filesize)
				_, err := io.ReadFull(reader, fileData)
				if err != nil {
					log.Printf("Error reading file: %v", err)
					continue
				}

				err = os.WriteFile(filename, fileData, 0644)
				if err != nil {
					log.Printf("Error saving file: %v", err)
				} else {
					fmt.Printf("\n%s (%s) отправил файл: %s\n", sender, timestamp, filename)
				}
			} else {
				fmt.Print("\r" + header)
			}
		}
	}()

	consoleScanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter message to send (type 'exit' to quit or <filename> to send a file):")
	for consoleScanner.Scan() {
		text := consoleScanner.Text()
		if strings.ToLower(text) == "exit" {
			break
		}

		if strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">") {
			filename := strings.Trim(text, "<>")
			sendFile(conn, filename, nickname)
			continue
		}

		if _, err := conn.Write([]byte(text + "\n")); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}

	if err := consoleScanner.Err(); err != nil {
		return fmt.Errorf("error reading from console %v", err.Error())
	}
	return nil
}

func sendFile(conn net.Conn, path string, nickname string) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	filename := fileInfo.Name()
	filesize := fileInfo.Size()
	timestamp := time.Now().Format("15:04:05")

	meta := fmt.Sprintf("FILE:%s|%s|%s|%d\n", nickname, timestamp, filename, filesize)
	_, err = conn.Write([]byte(meta))
	if err != nil {
		log.Printf("Error sending metadata: %v", err)
		return
	}

	if _, err := io.Copy(conn, file); err != nil {
		log.Printf("Error sending file data: %v", err)
	}
}
