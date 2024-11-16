package main

import (
	"net/http"
	"stockmark/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

func OnGetPortfolio(c echo.Context) error {
	permit := model.Permit(c.QueryParam("permit"))

	presentation, err := permit.PresentUser()
	if err != nil {
		return failure(c, err)
	}

	return c.JSON(http.StatusOK, object{
		"success": true,
		"data":    presentation,
	})
}

func OnGetDeposit(c echo.Context) error {
	permit := model.Permit(c.QueryParam("permit"))
	amount := c.QueryParam("amount")

	i64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return c.JSON(http.StatusOK, object{
			"success": false,
			"message": "invalid arguments",
		})
	}

	err = permit.Deposit(int(i64))
	if err != nil {
		return failure(c, err)
	}

	presentation, _ := permit.PresentUser()
	return c.JSON(http.StatusOK, object{
		"success": true,
		"data":    presentation,
	})
}
