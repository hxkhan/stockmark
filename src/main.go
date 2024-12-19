package main

import (
	"log"
	"net/http"
	"os"
	"stockmark/model"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// go build -o ./main.exe ./src && ./main serve
func main() {
	app := pocketbase.New()
	model.Initialize(app)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		se.Router.POST("/register", OnPostRegister)
		se.Router.POST("/login", OnPostLogin)

		se.Router.GET("/portfolio", OnGetPortfolio)
		se.Router.GET("/deposit", OnGetDeposit)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func failure(e *core.RequestEvent, err error) error {
	return e.JSON(http.StatusOK, map[string]any{
		"success": false,
		"message": err.Error(),
	})
}
