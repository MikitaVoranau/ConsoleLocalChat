package main

import (
	"CN_lab2/internal/client"
	"CN_lab2/internal/server"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting program...")
	fmt.Println("Choose mode to run program: ")
	fmt.Println("1. Starting server")
	fmt.Println("2. Starting client")
	fmt.Println("3.Press Ctrl+C to exit")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	switch choice {
	case "1":
		if err := server.StartServer(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
			fmt.Println("Press Enter to exit...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}
	case "2":
		if err := client.StartClient(); err != nil {
			fmt.Println("Error starting client:", err)
		}
	default:
		fmt.Println("Invalid choice")
	}

}
