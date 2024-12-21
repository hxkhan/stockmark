package model

import (
	"log"

	"github.com/pocketbase/dbx"
)

type User struct {
	Id             string  `db:"id"`
	Password       string  `db:"password" json:"-"`
	TokenKey       string  `db:"tokenKey"`
	Email          string  `db:"email"`
	FirstName      string  `db:"firstName"`
	LastName       string  `db:"lastName"`
	CurrentBalance float64 `db:"currentBalance"`
	TotalDeposited float64 `db:"totalDeposited"`
}

type Ticker struct {
	Id     string `db:"id"`
	Ticker string `db:"ticker"`
}

type Asset struct {
	Ticker      string  `json:"ticker"`
	Amount      int     `json:"amount"`
	AverageCost float64 `json:"averageCost"`
}

type Transaction struct {
	Id        string  `db:"id" json:"id"`
	Action    string  `db:"action" json:"action"`
	Ticker    string  `db:"ticker" json:"ticker"`
	Amount    int     `db:"amount" json:"amount"`
	UnitPrice float64 `db:"unitPrice" json:"unitPrice"`
	Time      string  `db:"time" json:"time"`
}

func (u User) FetchTransactions() []Transaction {
	transactions := []Transaction{}
	pb.DB().NewQuery(`
		SELECT transactions.id AS id, action, tickers.ticker AS ticker, amount, unitPrice, time
		FROM transactions
		LEFT JOIN tickers ON transactions.ticker = tickers.id
		WHERE transactions.actor = {:o}
		ORDER BY time DESC
	`).Bind(dbx.Params{
		"o": u.Id,
	}).All(&transactions)
	return transactions
}

func (user User) CalculateAssets() (assets []Asset) {
	allAssets := map[string]struct {
		amount int
		avg    float64
	}{}

	transactions := user.FetchTransactions()
	for i := len(transactions) - 1; i >= 0; i-- {
		transact := transactions[i]
		asset := allAssets[transact.Ticker]

		switch transact.Action {
		case "buy":
			// weighted average formula
			cost := asset.avg*float64(asset.amount) + transact.UnitPrice*float64(transact.Amount)
			quantity := asset.amount + transact.Amount
			asset.avg = cost / float64(quantity)
			asset.amount += transact.Amount
		case "sell":
			asset.amount -= transact.Amount
		default:
			log.Printf("CalculateAssets() -> unknown action '%v' \n", transact.Action)
		}

		allAssets[transact.Ticker] = asset
	}

	for tick, a := range allAssets {
		if a.amount > 0 {
			assets = append(assets, Asset{Ticker: tick, Amount: a.amount, AverageCost: a.avg})
		}
	}

	return assets
}
