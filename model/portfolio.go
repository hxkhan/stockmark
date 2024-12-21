package model

import (
	"math"
)

type PortfolioPage struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`

	CurrentBalance    float64 `json:"currentBalance"`
	DepositedTillDate float64 `json:"depositedTillDate"`

	TotalSpentOnActiveAssets float64 `json:"totalSpentOnActiveAssets"`
	ActiveAssetsValue        float64 `json:"activeAssetsValue"`
	ActiveAssetsValuePc      float64 `json:"activeAssetsValuePc"`
	ActiveAssetsProfit       float64 `json:"activeAssetsProfit"`

	MostProfitable       string  `json:"mostProfitable"`
	MostProfitableAmount float64 `json:"mostProfitableAmount"`
	MostProfitablePc     float64 `json:"mostProfitablePc"`

	LeastProfitable       string  `json:"leastProfitable"`
	LeastProfitableAmount float64 `json:"leastProfitableAmount"`
	LeastProfitablePc     float64 `json:"leastProfitablePc"`

	MostProfitableToday       string  `json:"mostProfitableToday"`
	MostProfitableTodayAmount float64 `json:"mostProfitableTodayAmount"`
	MostProfitableTodayPc     float64 `json:"mostProfitableTodayPc"`

	LeastProfitableToday       string  `json:"leastProfitableToday"`
	LeastProfitableTodayAmount float64 `json:"leastProfitableTodayAmount"`
	LeastProfitableTodayPc     float64 `json:"leastProfitableTodayPc"`

	TotalChangeToday float64 `json:"totalChangeToday"`

	Assets []PortfolioPageAsset `json:"assets"`
}

type PortfolioPageAsset struct {
	Ticker       string  `json:"ticker"`
	Amount       int     `json:"amount"`
	TotalWorth   float64 `json:"totalWorth"`
	PricePerUnit float64 `json:"pricePerUnit"`
	ChangePc     float64 `json:"changePc"`
}

func (p Permit) RenderPortfolio() (PortfolioPage, error) {
	acc, err := p.exchange()
	if err != nil {
		return PortfolioPage{}, err
	}

	assets := acc.CalculateAssets()

	user := PortfolioPage{FirstName: acc.FirstName, LastName: acc.LastName, Email: acc.Email, CurrentBalance: acc.CurrentBalance, DepositedTillDate: acc.TotalDeposited}
	userHasAssets := len(assets) != 0

	// prepare min and max values
	if userHasAssets {
		user.MostProfitableAmount = -math.MaxFloat64
		user.LeastProfitableAmount = math.MaxFloat64
		user.MostProfitableTodayAmount = -math.MaxFloat64
		user.LeastProfitableTodayAmount = math.MaxFloat64
	}

	for _, asset := range assets {
		stock, _ := Inquire(asset.Ticker)
		boughtFor := asset.AverageCost * float64(asset.Amount)
		currentvalue := stock.LastPrice * float64(asset.Amount)
		yesterdayValue := currentvalue / (1 + stock.PercentChange)
		differenceSinceBought := currentvalue - boughtFor
		differenceSinceYesterday := currentvalue - yesterdayValue

		user.TotalSpentOnActiveAssets += boughtFor
		user.ActiveAssetsValue += currentvalue

		user.TotalChangeToday += differenceSinceYesterday

		if differenceSinceBought > user.MostProfitableAmount {
			user.MostProfitable = asset.Ticker
			user.MostProfitableAmount = differenceSinceBought
			user.MostProfitablePc = ((currentvalue / boughtFor) - 1) * 100
		}
		if differenceSinceYesterday > user.MostProfitableTodayAmount {
			user.MostProfitableToday = asset.Ticker
			user.MostProfitableTodayAmount = differenceSinceYesterday
			user.MostProfitableTodayPc = stock.PercentChange * 100
			//user.MostProfitableTodayPc = ((currentvalue / yesterdayValue) - 1) * 100
		}

		if differenceSinceBought < user.LeastProfitableAmount {
			user.LeastProfitable = asset.Ticker
			user.LeastProfitableAmount = differenceSinceBought
			user.LeastProfitablePc = ((currentvalue / boughtFor) - 1) * 100
		}
		if differenceSinceYesterday < user.LeastProfitableTodayAmount {
			user.LeastProfitableToday = asset.Ticker
			user.LeastProfitableTodayAmount = differenceSinceYesterday
			user.LeastProfitableTodayPc = stock.PercentChange * 100
			//user.LeastProfitableTodayPc = ((currentvalue / yesterdayValue) - 1) * 100
		}

		user.Assets = append(user.Assets, PortfolioPageAsset{asset.Ticker, asset.Amount, currentvalue, stock.LastPrice, stock.PercentChange * 100})
	}

	// mostly roundings
	if userHasAssets {
		user.ActiveAssetsValuePc = ((user.ActiveAssetsValue - user.TotalSpentOnActiveAssets) / math.Abs(user.TotalSpentOnActiveAssets)) * 100
		user.ActiveAssetsProfit = user.ActiveAssetsValue - user.TotalSpentOnActiveAssets

		user.TotalSpentOnActiveAssets = math.Round(user.TotalSpentOnActiveAssets*100) / 100
		user.ActiveAssetsValue = math.Round(user.ActiveAssetsValue*100) / 100
		user.ActiveAssetsValuePc = math.Round(user.ActiveAssetsValuePc*100) / 100
		user.ActiveAssetsProfit = math.Round(user.ActiveAssetsProfit*100) / 100

		user.MostProfitableAmount = math.Round(user.MostProfitableAmount*100) / 100
		user.MostProfitablePc = math.Round(user.MostProfitablePc*100) / 100

		user.LeastProfitableAmount = math.Round(user.LeastProfitableAmount*100) / 100
		user.LeastProfitablePc = math.Round(user.LeastProfitablePc*100) / 100

		user.MostProfitableTodayAmount = math.Round(user.MostProfitableTodayAmount*100) / 100
		user.MostProfitableTodayPc = math.Round(user.MostProfitableTodayPc*100) / 100

		user.LeastProfitableTodayAmount = math.Round(user.LeastProfitableTodayAmount*100) / 100
		user.LeastProfitableTodayPc = math.Round(user.LeastProfitableTodayPc*100) / 100

		user.TotalChangeToday = math.Round(user.TotalChangeToday*100) / 100
	}

	return user, nil
}
