package model

import (
	"log"
	"stockmark/db"
)

func (p Permit) Deposit(amount int) error {
	acc, err := p.exchange()
	if err != nil {
		return err
	}

	acc.Deposited += float64(amount)
	acc.Balance += float64(amount)

	// impossible to fail
	err = db.UpdateAccount(acc)
	if err != nil {
		log.Printf("func Deposit() -> Permit exchanged successfully but update failed for user '%v'", acc.Email)
	}

	return nil
}
