package model

import (
	"github.com/pocketbase/dbx"
)

type User struct {
	Id             string  `db:"id"`
	Password       string  `db:"password"`
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
	Id          string  `db:"id"`
	Ticker      string  `db:"ticker"`
	Amount      int     `db:"amount"`
	AverageCost float64 `db:"averageCost"`
}

func (u User) FetchAssets() []Asset {
	assets := []Asset{}
	pb.DB().NewQuery(`
		SELECT assets.id AS id, tickers.ticker AS ticker, amount, averageCost
		FROM assets
		LEFT JOIN users ON assets.owner = users.id
		LEFT JOIN tickers ON assets.ticker = tickers.id
		WHERE assets.owner = {:o}
	`).Bind(dbx.Params{
		"o": u.Id,
	}).All(&assets)
	return assets
}

func (u *User) Deposit(amount int) error {
	u.CurrentBalance += float64(amount)
	u.TotalDeposited += float64(amount)

	_, err := pb.DB().NewQuery(`
		UPDATE users 
		SET currentBalance = {:cb}, totalDeposited = {:td}
		WHERE users.id = {:id}
	`).Bind(dbx.Params{
		"cb": u.CurrentBalance,
		"td": u.TotalDeposited,
		"id": u.Id,
	}).Execute()
	return err
}
