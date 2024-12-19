package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
)

var logins = make(map[Permit]entry)

type Permit string

type entry struct {
	email        string
	lastActivity time.Time
}

func Register(fname, lname, email, password string) (Permit, error) {
	err := pb.DB().NewQuery("SELECT * FROM users WHERE email = {:e}").Bind(dbx.Params{
		"e": email,
	}).One(&User{})

	// no existing user with the specified email
	if err == sql.ErrNoRows {
		_, err = pb.DB().NewQuery("INSERT INTO users (email, firstName, lastName, password) VALUES ({:e}, {:f}, {:l}, {:p})").Bind(dbx.Params{
			"e": email,
			"f": fname,
			"l": lname,
			"p": password,
		}).Execute()

		if err != nil {
			return "", err
		}

		id := Permit(uuid.New().String())
		logins[id] = entry{
			email,
			time.Now(),
		}

		return id, nil
	}

	// no error means user with the specified email exists
	if err == nil {
		return "", ErrAccountExists
	}

	// unknown error
	return "", ErrAccountExists
}

func Login(email, password string) (Permit, error) {
	user := User{}
	err := pb.DB().NewQuery("SELECT * FROM users WHERE email = {:e}").Bind(dbx.Params{
		"e": email,
	}).One(&user)

	if err == sql.ErrNoRows {
		return "", ErrAccountNotExists
	}

	if user.Password != password {
		return "", ErrIncorrectPassword
	}

	id := Permit(uuid.New().String())
	logins[id] = entry{
		email,
		time.Now(),
	}

	return id, nil
}

func (p Permit) exchange() (User, error) {
	if p == "" {
		return User{}, ErrNoPermitProvided
	}

	entry, has := logins[p]
	if !has {
		return User{}, ErrNotLoggedIn
	}

	user := User{}
	err := pb.DB().NewQuery("SELECT * FROM users WHERE email = {:e}").Bind(dbx.Params{
		"e": entry.email,
	}).One(&user)

	if err == sql.ErrNoRows {
		return User{}, ErrAccountNotExists
	}

	return user, nil
}

func (p Permit) Logout() error {
	if p == "" {
		return ErrNoPermitProvided
	}

	delete(logins, p)
	return nil
}
