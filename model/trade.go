package model

type TradePage struct {
	CurrentBalance float64                  `json:"currentBalance"`
	Transactions   []Transaction            `json:"transactions"`
	Market         []StockPresentationBasic `json:"market"`
	Assets         []Asset                  `json:"assets"`
}

type TradePageAsset struct {
	Ticker string `json:"ticker"`
	Amount int    `json:"amount"`
}

func (p Permit) RenderTrade() (page TradePage, err error) {
	acc, err := p.exchange()
	if err != nil {
		return page, err
	}

	page.CurrentBalance = acc.CurrentBalance

	// add user transactions
	page.Transactions = acc.FetchTransactions()

	// add all market stocks
	for _, stock := range stocks {
		stock.PercentChange *= 100
		page.Market = append(page.Market, stock)
	}

	// add user assets
	page.Assets = acc.CalculateAssets()

	return page, nil
}
