package model

import (
	"time"

	"github.com/pocketbase/dbx"
)

func (p Permit) Deposit(amount int) error {
	user, err := p.exchange()
	if err != nil {
		return err
	}

	_, err = pb.DB().NewQuery(`
		UPDATE users 
		SET currentBalance = currentBalance + {:amount}, totalDeposited = totalDeposited + {:amount}
		WHERE users.id = {:id}
	`).Bind(dbx.Params{
		"amount": amount,
		"id":     user.Id,
	}).Execute()
	return err
}

func (p Permit) Buy(ticker string, amount int) error {
	user, err := p.exchange()
	if err != nil {
		return err
	}

	stock, exists := stocks[ticker]
	if !exists {
		return ErrUnknownTicker
	}

	cost := stock.LastPrice * float64(amount)

	if user.CurrentBalance < cost {
		return ErrInsufficientBalance
	}

	_, err = pb.DB().NewQuery("INSERT INTO transactions (actor, action, ticker, amount, unitPrice, time) VALUES ({:actor}, {:action}, {:ticker}, {:amount}, {:unitPrice}, {:time})").Bind(dbx.Params{
		"actor":     user.Id,
		"action":    "buy",
		"ticker":    stock.Id,
		"amount":    amount,
		"unitPrice": stock.LastPrice,
		"time":      time.Now().UTC().Format(time.DateTime),
	}).Execute()
	return err
}

func (p Permit) Sell(ticker string, amount int) error {
	user, err := p.exchange()
	if err != nil {
		return err
	}

	stock, exists := stocks[ticker]
	if !exists {
		return ErrUnknownTicker
	}

	have := 0
	for _, asset := range user.CalculateAssets() {
		if asset.Ticker == ticker {
			have = asset.Amount
		}
	}

	if have < amount {
		return ErrInsufficientStocks
	}

	_, err = pb.DB().NewQuery("INSERT INTO transactions (actor, action, ticker, amount, unitPrice, time) VALUES ({:actor}, {:action}, {:ticker}, {:amount}, {:unitPrice}, {:time})").Bind(dbx.Params{
		"actor":     user.Id,
		"action":    "sell",
		"ticker":    stock.Id,
		"amount":    amount,
		"unitPrice": stock.LastPrice,
		"time":      time.Now().UTC().Format(time.DateTime),
	}).Execute()
	return err
}
