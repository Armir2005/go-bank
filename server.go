package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	Encoder    *json.Encoder
	WriterLock sync.Mutex
}

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.Mutex
)

type Command struct {
	Action   string
	Args     []string
	ClientID int
}

type Response struct {
	Success bool
	Message string
	Data    interface{}
}

type Session struct {
	Account *Account
}

func startServer() {
	bank := NewBank()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Bank server running on :8080")
	var clientID int

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}

		clientID++
		go handleClient(conn, bank, clientID)
	}
}

func handleClient(conn net.Conn, bank *Bank, clientID int) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	session := &Session{}

	client := &Client{Encoder: encoder}
	clientsMutex.Lock()
	clients[clientID] = client
	clientsMutex.Unlock()

	fmt.Printf("Client %d connected\n", clientID)

	for {
		var cmd Command
		err := decoder.Decode(&cmd)
		if err != nil {
			log.Printf("Client %d disconnected", clientID)
			return
		}

		response := executeCommand(bank, cmd, session)
		client.WriterLock.Lock()
		encoder.Encode(response)
		client.WriterLock.Unlock()

		if strings.ToLower(cmd.Action) == "exit" {
			log.Printf("Client %d disconnected", clientID)
			return
		}
	}
}

func executeCommand(bank *Bank, cmd Command, session *Session) Response {
	switch cmd.Action {
	case "help":
		helpMsg := "Available commands:\n" +
			"  help                           - Show available commands\n" +
			"  create [Name] [Password]       - Create an account\n" +
			"  login [AccountID] [Password]   - Login to your account\n" +
			"  logout                         - Logout from your account\n" +
			"  deposit [Amount] [Note]        - Deposit money\n" +
			"  withdraw [Amount] [Note]       - Withdraw money\n" +
			"  transfer [toID] [Amount] [Note]- Transfer money\n" +
			"  balance                        - Show account balance\n" +
			"  transactions                   - Show transaction history\n" +
			"  message [Text]                 - Send a message to all clients\n" +
			"  exit                           - Exit the application"
		return Response{true, helpMsg, nil}
	case "create":
		if len(cmd.Args) < 2 {
			return Response{false, "Usage: create [Name] [Password]", nil}
		}
		acc := bank.CreateAccount(cmd.Args[0], cmd.Args[1])
		return Response{true, fmt.Sprintf("Account created with ID %d", acc.ID), acc.ID}

	case "login":
		if len(cmd.Args) < 2 {
			return Response{false, "Usage: login [AccountID] [Password]", nil}
		}
		id, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return Response{false, "Invalid account ID", nil}
		}
		acc, ok := bank.GetAccount(id)
		if !ok {
			return Response{false, "Account not found", nil}
		}
		if acc.Password != cmd.Args[1] {
			return Response{false, "Invalid password", nil}
		}
		session.Account = acc
		return Response{true, fmt.Sprintf("Logged in as %s", acc.Owner), nil}
	case "logout":
		if session.Account == nil {
			return Response{false, "Not logged in", nil}
		}
		session.Account = nil
		return Response{true, "Logged out successfully", nil}
	case "deposit":
		if session.Account == nil {
			return Response{false, "Please login first", nil}
		}
		if len(cmd.Args) < 2 {
			return Response{false, "Usage: deposit [Amount] [Note]", nil}
		}
		amount, err := strconv.ParseFloat(cmd.Args[0], 64)
		if err != nil {
			return Response{false, "Invalid amount", nil}
		}
		note := strings.Join(cmd.Args[1:], " ")
		if err := session.Account.Deposit(amount, note); err != nil {
			return Response{false, "Error depositing: " + err.Error(), nil}
		}
		return Response{true, "Deposit successful", nil}
	case "withdraw":
		if session.Account == nil {
			return Response{false, "Please login first", nil}
		}
		if len(cmd.Args) < 2 {
			return Response{false, "Usage: withdraw [Amount] [Note]", nil}
		}
		amount, err := strconv.ParseFloat(cmd.Args[0], 64)
		if err != nil {
			return Response{false, "Invalid amount", nil}
		}
		note := strings.Join(cmd.Args[1:], " ")
		if err := session.Account.Withdraw(amount, note); err != nil {
			return Response{false, "Error withdrawing: " + err.Error(), nil}
		}
		return Response{true, "Withdrawal successful", nil}
	case "transfer":
		if session.Account == nil {
			return Response{false, "Please login first", nil}
		}
		if len(cmd.Args) < 3 {
			return Response{false, "Usage: transfer [toID] [Amount] [Note]", nil}
		}
		toID, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return Response{false, "Invalid account ID", nil}
		}
		amount, err := strconv.ParseFloat(cmd.Args[1], 64)
		if err != nil {
			return Response{false, "Invalid amount", nil}
		}
		note := strings.Join(cmd.Args[2:], " ")
		if err := bank.Transfer(session.Account.ID, toID, amount, note); err != nil {
			return Response{false, "Error transferring: " + err.Error(), nil}
		}
		return Response{true, "Transfer successful", nil}
	case "balance":
		if session.Account == nil {
			return Response{false, "Please login first", nil}
		}
		msg := fmt.Sprintf("Your balance is %.2f", session.Account.Balance)
		return Response{true, msg, session.Account.Balance}
	case "transactions":
		if session.Account == nil {
			return Response{false, "Please login first", nil}
		}
		var result strings.Builder
		result.WriteString("Transaction for " + session.Account.Owner + ":\n")
		for _, t := range session.Account.Transactions {
			result.WriteString(fmt.Sprintf(" - %s: %.2f at %s (%s)\n", t.Type, t.Amount, t.Timestamp.Format("2006-01-02 15:04:05"), t.Note))
		}
		result.WriteString(fmt.Sprintf("Balance: %.2f", session.Account.Balance))
		return Response{true, result.String(), nil}
	case "message":
		if session.Account == nil {
			return Response{false, "Please login first", nil}
		}
		if len(cmd.Args) < 1 {
			return Response{false, "Usage: message [Text]", nil}
		}
		msgText := strings.Join(cmd.Args, " ")
		fullMsg := fmt.Sprintf("%s: %s", session.Account.Owner, msgText)
		broadcastMessage(Response{true, fullMsg, nil})
		return Response{true, "Message sent", nil}
	case "exit":
		return Response{true, "Goodbye!", nil}
	default:
		return Response{false, "Unknown command", nil}
	}
}

func broadcastMessage(response Response) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for _, client := range clients {
		client.WriterLock.Lock()
		client.Encoder.Encode(response)
		client.WriterLock.Unlock()
	}
}
