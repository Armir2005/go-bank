# go-bank

go-bank is a simple interactive banking simulation written in Go. It allows users to create accounts, deposit, withdraw, transfer money, and view account balances and transaction histories through a command-line interface.

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/Armir2005/go-bank.git
   cd go-bank
   ```

2. **Build the project:**

    You can run the program without building by using:

    ```sh
    go run *.go
    ```

    Or build an executable:

    ```sh
    go build -o go-bank
    ```

## Usage

1. **Start the program from your terminal:**

    ```sh
    ./go-bank
    ```

2. **When the program starts, you will see a prompt:**

    ```sh
    Welcome to the go-bank system!
    Type 'help', to show the available commands.
    > 
    ```

### Available Commands

- **help**
  Displays a list of available commands.
- **create [Name] [Password]**
  Creates a new account for the specified owner.
  _Example:_
  `create Alice SuperSecurePassword`
- **deposit [Amount] [Note]**
  Deposits the specified amount of money into the account. A note is added to the transaction.
  _Example:_
  `deposit 1000 Initial deposit`
- **withdraw [Amount] [Note]**
  Withdrwas the specified amount of money from the account with a note.
  _Example:_
  `deposit 200 Rent payment`
- **transfer [toID] [Amount] [Note]**
  Transfers money from one account to another.
  _Example:_
  `transfer 2 300 Payment for goods`
- **balance**
  Shows the current balance of the specified account.
  _Example:_
  `balance`
- **transactions**
  Prints the transaction history for the account.
  _Example:_
  `transactions`
- **clear**
  Clears the terminal screen.
- **exit**
  Exits the program.

## License

This project is licensed under the MIT License.