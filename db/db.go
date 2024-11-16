package db

import (
	"encoding/json"
	"os"
	"strings"
)

const accountsFileName = "./db/accounts.json"
const tickersFileName = "./db/tickers.json"

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

func LoadAccount(email string) (acc Account, exists bool) {
	email = strings.ToLower(email)
	accounts := fetchUsers()

	acc, exists = accounts[email]
	acc.Email = email
	return acc, exists
}

func SaveAccount(acc Account) {
	accounts := fetchUsers()
	accounts[strings.ToLower(acc.Email)] = acc
	saveUsers(accounts)
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
