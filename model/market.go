package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"stockmark/db"
	"strings"
)

var stocks = make(map[string]StockPresentationBasic)

type StockPresentationBasic struct {
	Symbol        string  `json:"symbol"`
	SymbolName    string  `json:"symbolName"`
	LastPrice     float64 `json:"lastPrice"`
	PriceChange   string  `json:"priceChange"`
	PercentChange float64 `json:"percentChange"`
	LowPrice      float64 `json:"lowPrice"`
	OpenPrice     float64 `json:"openPrice"`
	HighPrice     float64 `json:"highPrice"`
	Volume        float64 `json:"volume"`
}

func init() {
	supportedTickers := db.GetSupportedTickers()

	for _, v := range supportedTickers {
		go fetchStock(v)
	}
}

func fetchStock(ticker string) {
	req, _ := http.NewRequest("GET", "https://realstonks.p.rapidapi.com/stocks/"+strings.ToUpper(ticker), nil)

	req.Header.Add("x-rapidapi-key", "bffff99a21msh3201484be4f30c1p1ab348jsn39aebb9ede1c")
	req.Header.Add("x-rapidapi-host", "realstonks.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("for ticker '%v', fetch could not be performed\n", err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("for ticker '%v', response body could not be read\n", err.Error())
		return
	}

	var s StockPresentationBasic
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Printf("for ticker '%v', error: %v\n", ticker, err.Error())
		return
	}

	// unsafe for concurrent access but who cares
	stocks[ticker] = s
	fmt.Printf("Retrieved %v\n", ticker)
}

func Inquire(ticker string) (StockPresentationBasic, error) {
	if ticker == "" {
		return StockPresentationBasic{}, ErrInvalidArguments
	}

	s, has := stocks[ticker]

	if !has {
		return StockPresentationBasic{}, ErrUnknownTicker
	}

	return s, nil
}
