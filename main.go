package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/atm_demo"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func register(username string, pin int) {
	_, err := db.Exec("INSERT INTO accounts (username, pin, balance) VALUES (?, ?, 0)", username, pin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registration successful!")
}

func login(username string, pin int) (int, bool) {
	var id int
	err := db.QueryRow("SELECT id FROM accounts WHERE username=? AND pin=?", username, pin).Scan(&id)
	if err != nil {
		fmt.Println("Login failed!")
		return 0, false
	}
	fmt.Println("Login successful!")
	return id, true
}

func checkBalance(id int) {
	var balance float64
	err := db.QueryRow("SELECT balance FROM accounts WHERE id=?", id).Scan(&balance)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Your current balance: %.2f\n", balance)
}

func deposit(id int, amount float64) {
	_, err := db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deposit successful!")
}

func withdraw(id int, amount float64) {
	var balance float64
	err := db.QueryRow("SELECT balance FROM accounts WHERE id=?", id).Scan(&balance)
	if err != nil {
		log.Fatal(err)
	}
	if balance < amount {
		fmt.Println("Insufficient balance!")
		return
	}
	_, err = db.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Withdrawal successful!")
}

func deleteAccount(id int) {
	_, err := db.Exec("DELETE FROM accounts WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Account deleted successfully!")
}

func main() {
	initDB()
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Usage: register|login [parameters]")
		return
	}

	switch os.Args[1] {
	case "register":
		if len(os.Args) != 4 {
			fmt.Println("Usage: register <username> <pin>")
			return
		}
		username := os.Args[2]
		var pin int
		fmt.Sscanf(os.Args[3], "%d", &pin)
		register(username, pin)
	case "login":
		if len(os.Args) != 4 {
			fmt.Println("Usage: login <username> <pin>")
			return
		}
		username := os.Args[2]
		var pin int
		fmt.Sscanf(os.Args[3], "%d", &pin)
		id, success := login(username, pin)
		if !success {
			return
		}
		// logged in
		fmt.Println("Select an action: check|deposit|withdraw|delete")
		var action string
		fmt.Scan(&action)

		switch action {
		case "check":
			checkBalance(id)
		case "deposit":
			var amount float64
			fmt.Print("Enter deposit amount: ")
			fmt.Scan(&amount)
			deposit(id, amount)
		case "withdraw":
			var amount float64
			fmt.Print("Enter withdrawal amount: ")
			fmt.Scan(&amount)
			withdraw(id, amount)
		case "delete":
			deleteAccount(id)
		default:
			fmt.Println("Unknown action!")
		}
	default:
		fmt.Println("Unknown command!")
	}
}
