package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	bank := NewBank()
	scanner := bufio.NewScanner(os.Stdin)

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
			fmt.Println("  create [Name]                            - Create an account")
			fmt.Println("  deposit [AccountID] [Amount] [Note]      - Deposit money")
			fmt.Println("  withdraw [AccountID] [Amount] [Note]     - Withdraw money")
			fmt.Println("  transfer [fromID] [toID] [Amount] [Note] - Transfer monye")
			fmt.Println("  balance [AccountID]                      - Show account balance")
			fmt.Println("  transactions [AccountID]                 - Show transactions")
			fmt.Println("  clear                                    - Clear the terminal")
			fmt.Println("  exit                                     - Exit the game")
		case "create":
			if len(parts) < 2 {
				fmt.Println("Usage: create [Name]")
				continue
			}
			owner := parts[1]
			acc := bank.CreateAccount(owner)
			fmt.Printf("Account created for %s with the ID %d\n", owner, acc.ID)
		case "deposit":
			if len(parts) < 4 {
				fmt.Println("Usage: deposit [AccountID] [Amount] [Note]")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid accountID")
				continue
			}
			amount, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			note := strings.Join(parts[3:], " ")
			acc, ok := bank.GetAccount(id)
			if !ok {
				fmt.Println("Account not found")
				continue
			}
			if err := acc.Deposit(amount, note); err != nil {
				fmt.Println("Something went wrong depositing:", err)
			} else {
				fmt.Println("Deposit went sucessfully!")
			}
		case "withdraw":
			if len(parts) < 4 {
				fmt.Println("Usage: withdraw [AccountID] [Amount] [Note]")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid accountID")
				continue
			}
			amount, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			note := strings.Join(parts[3:], " ")
			acc, ok := bank.GetAccount(id)
			if !ok {
				fmt.Println("Account not fount")
				continue
			}
			if err := acc.Withdraw(amount, note); err != nil {
				fmt.Println("Error while withdrawing:", err)
			} else {
				fmt.Println("Withdraw went sucessfully!")
			}
		case "transfer":
			if len(parts) < 5 {
				fmt.Println("Usage: transfer [fromID] [toID] [Amount] [Note]")
				continue
			}
			fromID, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid fromID")
				continue
			}
			toID, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid toID")
				continue
			}
			amount, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			note := strings.Join(parts[4:], " ")
			if err := bank.Transfer(fromID, toID, amount, note); err != nil {
				fmt.Println("Error while transfering:", err)
			} else {
				fmt.Println("Transfer went sucessfully!")
			}
		case "balance":
			if len(parts) < 2 {
				fmt.Println("Usage: balance [AccountID]")
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
			acc.PrintBalance()
		case "transactions":
			if len(parts) < 2 {
				fmt.Println("Usage: transactions [AccountID]")
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
			acc.PrintTransactions()
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
