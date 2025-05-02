package additonal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetIPandPorts() (string, int, error) {
	reader := bufio.NewReader(os.Stdin)
	var ip string
	fmt.Println("Enter IP-address (by default localhost(`127.0.0.1`)) : ")
	ip, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, fmt.Errorf("error reading IP address: %v", err)
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		fmt.Println("The server is running on localhost(IP-address: 127.0.0.1)")
		ip = "127.0.0.1"
	}
	fmt.Println("Enter Port (by default `8080`) : ")
	portStr, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, fmt.Errorf("error reading port: %v", err)
	}
	portStr = strings.TrimSpace(portStr)
	port := 8080
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			fmt.Println("Invalid port number")
			fmt.Println("The server is running on Port: 8080")
			port = 8080
		}
	}
	return ip, port, nil
}
