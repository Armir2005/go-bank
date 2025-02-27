package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		startSinglePlayer()
		return
	}

	switch os.Args[1] {
	case "server":
		startServer()
	case "client":
		if len(os.Args) < 3 {
			fmt.Println("Usage: client [server-address]")
			return
		}
		startClient(os.Args[2])
	default:
		startSinglePlayer()
	}
}

func startSinglePlayer() {
	bank := NewBank()
	scanner := bufio.NewScanner(os.Stdin)
	var currentAccount *Account

	fmt.Println("Welcome to the go-bank system!")
	fmt.Println("Type 'help', to show the available commands.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  create [Name] [Password]         - Create an account")
			fmt.Println("  login [AccountID] [Password]     - Login to an account")
			fmt.Println("  logout                           - Logout from an account")
			fmt.Println("  deposit [Amount] [Note]          - Deposit money")
			fmt.Println("  withdraw [Amount] [Note]         - Withdraw money")
			fmt.Println("  transfer [toID] [Amount] [Note]  - Transfer monye")
			fmt.Println("  balance                          - Show account balance")
			fmt.Println("  transactions                     - Show transactions")
			fmt.Println("  clear                            - Clear the terminal")
			fmt.Println("  exit                             - Exit the game")
		case "create":
			if len(parts) < 3 {
				fmt.Println("Usage: create [Name] [Password]")
				continue
			}
			owner := parts[1]
			password := parts[2]
			acc := bank.CreateAccount(owner, password)
			fmt.Printf("Account created for %s with the ID %d\n", owner, acc.ID)
		case "login":
			if len(parts) < 3 {
				fmt.Println("Usage: login [AccountID] [Password]")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid accountID")
				continue
			}
			acc, ok := bank.GetAccount(id)
			if !ok {
				fmt.Println("Account not found")
				continue
			}
			if acc.Password != parts[2] {
				fmt.Println("Invalid password")
				continue
			}
			currentAccount = acc
			fmt.Printf("Logged in as %s\n", acc.Owner)
		case "logout":
			if currentAccount == nil {
				fmt.Println("Not logged in")
				continue
			}
			fmt.Printf("Logged out from %s\n", currentAccount.Owner)
			currentAccount = nil
		case "deposit":
			if currentAccount == nil {
				fmt.Println("Please login first.")
				continue
			}
			if len(parts) < 3 {
				fmt.Println("Usage: deposit [Amount] [Note]")
				continue
			}
			amount, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			note := strings.Join(parts[2:], " ")
			if err := currentAccount.Deposit(amount, note); err != nil {
				fmt.Println("Something went wrong depositing:", err)
			} else {
				fmt.Println("Deposit went sucessfully!")
			}
		case "withdraw":
			if currentAccount == nil {
				fmt.Println("Please login first.")
				continue
			}
			if len(parts) < 3 {
				fmt.Println("Usage: withdraw [Amount] [Note]")
				continue
			}
			amount, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			note := strings.Join(parts[2:], " ")
			if err := currentAccount.Withdraw(amount, note); err != nil {
				fmt.Println("Error while withdrawing:", err)
			} else {
				fmt.Println("Withdraw went sucessfully!")
			}
		case "transfer":
			if currentAccount == nil {
				fmt.Println("Please login first.")
				continue
			}
			if len(parts) < 4 {
				fmt.Println("Usage: transfer [toID] [Amount] [Note]")
				continue
			}
			toID, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid toID")
				continue
			}
			amount, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			note := strings.Join(parts[3:], " ")
			if err := bank.Transfer(currentAccount.ID, toID, amount, note); err != nil {
				fmt.Println("Error while transfering:", err)
			} else {
				fmt.Println("Transfer went sucessfully!")
			}
		case "balance":
			if currentAccount == nil {
				fmt.Println("Please login first.")
				continue
			}
			currentAccount.PrintBalance()
		case "transactions":
			if currentAccount == nil {
				fmt.Println("Please login first.")
				continue
			}
			currentAccount.PrintTransactions()
		case "clear":
			fmt.Print("\033[H\033[2J")
			continue
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command. Type 'help' to get help.")
		}
	}
}
