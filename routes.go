package main

import (
	"net/http"
	"stockmark/model"
	"strconv"

	"github.com/pocketbase/pocketbase/core"
)

func OnPostLogin(e *core.RequestEvent) error {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := e.BindBody(&body)
	if err != nil {
		return failure(e, err)
	}

	permit, err := model.Login(body.Email, body.Password)
	if err != nil {
		return failure(e, err)
	}

	page, err := permit.RenderPortfolio()
	if err != nil {
		return failure(e, err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"permit":  permit,
		"data":    page,
	})
}

func OnPostRegister(e *core.RequestEvent) error {
	body := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{}

	err := e.BindBody(&body)
	if err != nil {
		return failure(e, err)
	}

	permit, err := model.Register(body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		return failure(e, err)
	}

	presentation, err := permit.RenderPortfolio()
	if err != nil {
		return failure(e, err)
	}
	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"permit":  permit,
		"data":    presentation,
	})
}

func OnGetPortfolio(e *core.RequestEvent) error {
	permit := model.Permit(e.Request.URL.Query().Get("permit"))

	page, err := permit.RenderPortfolio()
	if err != nil {
		return failure(e, err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    page,
	})
}

func OnGetTrade(e *core.RequestEvent) error {
	permit := model.Permit(e.Request.URL.Query().Get("permit"))

	page, err := permit.RenderTrade()
	if err != nil {
		return failure(e, err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    page,
	})
}

func OnGetTransact(e *core.RequestEvent) error {
	query := e.Request.URL.Query()

	permit := model.Permit(query.Get("permit"))
	action := query.Get("action")
	ticker := query.Get("ticker")
	amount := query.Get("amount")

	i64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return failure(e, model.ErrInvalidArguments)
	}

	switch action {
	case "buy":
		err = permit.Buy(ticker, int(i64))
	case "sell":
		err = permit.Sell(ticker, int(i64))
	default:
		err = model.ErrInvalidArguments
	}

	if err != nil {
		return failure(e, err)
	}

	page, _ := permit.RenderTrade()
	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    page,
	})
}

func OnGetDeposit(e *core.RequestEvent) error {
	permit := model.Permit(e.Request.URL.Query().Get("permit"))
	amount := e.Request.URL.Query().Get("amount")

	i64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return failure(e, model.ErrInvalidArguments)
	}

	err = permit.Deposit(int(i64))
	if err != nil {
		return failure(e, err)
	}
	return e.JSON(http.StatusOK, map[string]any{"success": true})
}
