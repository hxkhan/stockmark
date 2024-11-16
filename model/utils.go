package model

import (
	"stockmark/db"
)

func (p Permit) Deposit(amount int) error {
	acc, err := p.exchange()
	if err != nil {
		return err
	}

	acc.Deposited += float64(amount)
	acc.Balance += float64(amount)

	db.SaveAccount(acc)
	return nil
}
