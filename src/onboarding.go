package main

import (
	"net/http"
	"stockmark/model"

	"github.com/labstack/echo/v4"
)

type BodyOfLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BodyOfRegister struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func OnPostLogin(c echo.Context) error {
	body, err := echoExtractBodyInto[BodyOfLogin](c)
	if err != nil {
		return err
	}

	permit, err := model.Login(body.Email, body.Password)
	if err != nil {
		return failure(c, err)
	}

	presentation, err := permit.PresentUser()
	if err != nil {
		return failure(c, err)
	}

	return c.JSON(http.StatusOK, object{
		"success": true,
		"permit":  permit,
		"data":    presentation,
	})
}

func OnPostRegister(c echo.Context) error {
	body, err := echoExtractBodyInto[BodyOfRegister](c)
	if err != nil {
		return err
	}

	permit, err := model.Register(body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		return failure(c, err)
	}

	presentation, err := permit.PresentUser()
	if err != nil {
		return failure(c, err)
	}
	return c.JSON(http.StatusOK, object{
		"success": true,
		"permit":  permit,
		"data":    presentation,
	})
}
