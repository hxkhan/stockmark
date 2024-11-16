package model

import (
	"stockmark/db"
	"time"

	"github.com/google/uuid"
)

var logins = make(map[Permit]entry)

type Permit string

type entry struct {
	email      string
	lastUpdate time.Time
}

func Register(fname, lname, email, password string) (Permit, error) {
	acc, err := db.CreateAccount(fname, lname, email, password)
	if err != nil {
		return "", err
	}

	id := Permit(uuid.New().String())
	logins[id] = entry{
		acc.Email,
		time.Now(),
	}

	return id, nil
}

func Login(email, password string) (Permit, error) {
	acc, err := db.GetAccount(email)
	if err != nil {
		return "", err
	}

	if acc.Password != password {
		return "", ErrIncorrectPassword
	}

	id := Permit(uuid.New().String())
	logins[id] = entry{
		acc.Email,
		time.Now(),
	}

	return id, nil
}

func (p Permit) exchange() (db.Account, error) {
	if p == "" {
		return db.Account{}, ErrNoPermitProvided
	}

	if entry, has := logins[p]; has {
		// also check how old the entry is
		acc, _ := db.GetAccount(entry.email)
		return acc, nil
	}
	return db.Account{}, ErrNotLoggedIn
}

func (p Permit) Logout() error {
	if p == "" {
		return ErrNoPermitProvided
	}

	delete(logins, p)
	return nil
}
