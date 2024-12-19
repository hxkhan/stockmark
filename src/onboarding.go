package main

import (
	"net/http"
	"stockmark/model"

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
