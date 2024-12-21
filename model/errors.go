package model

import "errors"

var ErrAccountExists = errors.New("account exists")
var ErrAccountNotExists = errors.New("account not exists")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrNoPermitProvided = errors.New("no permit provided")
var ErrNotLoggedIn = errors.New("not logged in")
var ErrUnknownTicker = errors.New("unknown ticker")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrInsufficientStocks = errors.New("insufficient stocks")
var ErrInvalidArguments = errors.New("invalid arguments")
