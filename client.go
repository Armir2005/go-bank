package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func startClient(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to bank server")

	go readResponses(conn)
	writeCommands(conn)
}

func readResponses(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	for {
		var res Response
		err := decoder.Decode(&res)
		if err != nil {
			fmt.Println("Disconnected from server")
			os.Exit(0)
		}
		fmt.Println("\nServer:", res.Message)
		fmt.Print("> ")
	}
}

func writeCommands(conn net.Conn) {
	encoder := json.NewEncoder(conn)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Fields(text)
		if len(parts) == 0 {
			continue
		}

		cmd := Command{
			Action: parts[0],
			Args:   parts[1:],
		}

		err := encoder.Encode(cmd)
		if err != nil {
			fmt.Println("Error encoding command:", err)
			continue
		}
		fmt.Print("> ")
	}
}
