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
	acc, exists := db.LoadAccount(email)
	if exists {
		return "", ErrAccountExists
	}

	db.SaveAccount(db.Account{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})

	id := Permit(uuid.New().String())
	logins[id] = entry{
		acc.Email,
		time.Now(),
	}

	return id, nil
}

func Login(email, password string) (Permit, error) {
	acc, exists := db.LoadAccount(email)
	if !exists {
		return "", ErrAccountNotExists
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

	entry, has := logins[p]
	if !has {
		return db.Account{}, ErrNotLoggedIn
	}

	acc, exists := db.LoadAccount(entry.email)
	if !exists {
		return db.Account{}, ErrAccountNotExists
	}

	return acc, nil

}

func (p Permit) Logout() error {
	if p == "" {
		return ErrNoPermitProvided
	}

	delete(logins, p)
	return nil
}
