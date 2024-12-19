package main

import (
	"net/http"
	"stockmark/model"
	"strconv"

	"github.com/pocketbase/pocketbase/core"
)

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

func OnGetDeposit(e *core.RequestEvent) error {
	permit := model.Permit(e.Request.URL.Query().Get("permit"))
	amount := e.Request.URL.Query().Get("amount")

	i64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return e.JSON(http.StatusOK, map[string]any{
			"success": false,
			"message": "invalid arguments",
		})
	}

	err = permit.Deposit(int(i64))
	if err != nil {
		return failure(e, err)
	}

	page, _ := permit.RenderPortfolio()
	return e.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    page,
	})
}
