package model

import (
	"log"
	"time"

	"github.com/pocketbase/pocketbase"
)

var pb *pocketbase.PocketBase

func Initialize(app *pocketbase.PocketBase) {
	pb = app

	// fetch and update tickers
	go func() {
		// for some reason the db connection takes a moment to be established
		time.Sleep(time.Second)
		updateTickers()
	}()
}

func updateTickers() {
	tickers := []Ticker{}
	err := pb.DB().NewQuery("SELECT * FROM tickers").All(&tickers)
	if err != nil {
		log.Printf("updateTickers() -> '%v' \n", err.Error())
		return
	}

	if len(tickers) == 0 {
		log.Println("No tickers in database!")
	}

	for _, v := range tickers {
		go fetchStock(v.Ticker)
	}
}
