package main

import (
	"net/http"
	"stockmark/model"

	"github.com/labstack/echo/v4"
)

type object map[string]any

// go build -o bin/main.exe ./src && bin\main
func main() {
	e := echo.New()
	e.Static("/", "public")

	e.POST("/login", OnPostLogin)
	e.POST("/register", OnPostRegister)
	e.GET("/logout", func(c echo.Context) error {
		permit := model.Permit(c.QueryParam("permit"))

		permit.Logout()
		return c.JSON(http.StatusOK, object{"success": true})
	})

	e.GET("/portfolio", OnGetPortfolio)
	e.GET("/deposit", OnGetDeposit)
	e.GET("/inquire", OnGetInquire)

	e.Logger.Fatal(e.Start(":8080"))
}

func echoExtractBodyInto[T any](c echo.Context) (T, error) {
	var value T
	err := c.Bind(&value)
	if err != nil {
		return value, c.JSON(http.StatusOK, object{
			"success": false,
			"message": "invalid request body",
		})
	}
	return value, nil
}

func failure(c echo.Context, err error) error {
	return c.JSON(http.StatusOK, object{
		"success": false,
		"message": err.Error(),
	})
}
