package db

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

const accountsFileName = "./db/accounts.json"
const tickersFileName = "./db/tickers.json"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

type Account struct {
	FirstName string
	LastName  string
	Email     string `json:"-"`
	Password  string
	Balance   float64
	Deposited float64
	Assets    map[string]struct {
		Amount   int
		BuyPrice float64
	}
	History []struct {
		Action    string
		Ticker    string
		Amount    int
		UnitPrice float64
		Time      string
	}
}

func GetAccount(email string) (Account, error) {
	email = strings.ToLower(email)
	accounts := fetchUsers()

	if acc, exists := accounts[email]; exists {
		acc.Email = email
		return acc, nil
	}
	return Account{}, ErrUserNotFound
}

func CreateAccount(name, lname, email, password string) (Account, error) {
	email = strings.ToLower(email)
	accounts := fetchUsers()

	if acc, exists := accounts[email]; exists {
		return acc, ErrUserAlreadyExists
	}

	acc := Account{FirstName: name, LastName: lname, Email: email, Password: password}
	accounts[email] = acc
	saveUsers(accounts)
	return acc, nil
}

func UpdateAccount(acc Account) error {
	accounts := fetchUsers()

	if _, exists := accounts[acc.Email]; !exists {
		return ErrUserNotFound
	}

	accounts[acc.Email] = acc
	saveUsers(accounts)
	return nil
}

func GetSupportedTickers() []string {
	bytes, err := os.ReadFile(tickersFileName)
	if err != nil {
		panic(err)
	}

	var tickers []string
	err = json.Unmarshal(bytes, &tickers)
	if err != nil {
		panic(err)
	}

	return tickers
}

func fetchUsers() map[string]Account {
	bytes, err := os.ReadFile(accountsFileName)
	if err != nil {
		panic(err)
	}

	var accounts map[string]Account
	err = json.Unmarshal(bytes, &accounts)
	if err != nil {
		panic(err)
	}

	return accounts
}

func saveUsers(accounts map[string]Account) {
	bytes, err := json.Marshal(accounts)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(accountsFileName, bytes, 0666)
	if err != nil {
		panic(err)
	}
}
