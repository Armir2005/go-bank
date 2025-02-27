package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const spacer = "*******************************"

type Account struct {
	ID           int
	Owner        string
	Password     string
	Balance      float64
	Transactions []Transaction
	mu           sync.Mutex
}

type Transaction struct {
	Type      string
	Amount    float64
	Timestamp time.Time
	Note      string
}

type Bank struct {
	accounts map[int]*Account
	mu       sync.RWMutex
	nextID   int
}

func NewBank() *Bank {
	return &Bank{
		accounts: make(map[int]*Account),
		nextID:   1,
	}
}

func (b *Bank) CreateAccount(owner, password string) *Account {
	b.mu.Lock()
	defer b.mu.Unlock()

	account := &Account{
		ID:       b.nextID,
		Owner:    owner,
		Password: password,
		Balance:  0,
	}
	b.accounts[b.nextID] = account
	b.nextID++
	return account
}

func (b *Bank) GetAccount(id int) (*Account, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	acc, ok := b.accounts[id]
	return acc, ok
}

func (a *Account) Deposit(amount float64, note string) error {
	if amount <= 0 {
		return errors.New("ammount must be positive")
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Balance += amount
	a.Transactions = append(a.Transactions, Transaction{
		Type:      "deposit",
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      note,
	})
	return nil
}

func (a *Account) Withdraw(amount float64, note string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if a.Balance < amount {
		return errors.New("not enough funds")
	}
	a.Balance -= amount
	a.Transactions = append(a.Transactions, Transaction{
		Type:      "withdraw",
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      note,
	})
	return nil
}

func (b *Bank) Transfer(fromID, toID int, amount float64, note string) error {
	b.mu.RLock()
	from, okFrom := b.accounts[fromID]
	to, okTo := b.accounts[toID]
	b.mu.RUnlock()

	if !okFrom || !okTo {
		return errors.New("account not found")
	}

	if err := from.Withdraw(amount, note); err != nil {
		return err
	}
	if err := to.Deposit(amount, note); err != nil {
		return err
	}

	return nil
}

func (a *Account) PrintTransactions() {
	a.mu.Lock()
	defer a.mu.Unlock()

	fmt.Println(spacer)
	fmt.Println("Transactions:", a.Owner)
	for _, t := range a.Transactions {
		fmt.Printf(" - %s: %.2f at %s (%s)\n", t.Type, t.Amount, t.Timestamp.Format("2006-01-02 15:04:05"), t.Note)
	}
	fmt.Println("Balance:", a.Balance)
	fmt.Println(spacer)
}

func (a *Account) PrintBalance() {
	a.mu.Lock()
	defer a.mu.Unlock()

	fmt.Println(spacer)
	fmt.Println(a.Owner + ":")
	fmt.Println("Balance:", a.Balance)
	fmt.Println(spacer)
}
