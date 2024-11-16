package model

import (
	"errors"
	"stockmark/db"
)

var ErrUserAlreadyExists = db.ErrUserAlreadyExists
var ErrUserNotFound = db.ErrUserNotFound

var ErrIncorrectPassword = errors.New("incorrect password")
var ErrNoPermitProvided = errors.New("no permit provided")
var ErrNotLoggedIn = errors.New("not logged in")
var ErrUnknownTicker = errors.New("unknown ticker")
var ErrInvalidArguments = errors.New("invalid arguments")
