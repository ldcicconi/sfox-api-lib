package sfoxapi

import (
	"github.com/shopspring/decimal"
)

type SFOXBalance struct {
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
}

func (api *SFOXAPI) GetBalances() (balances []SFOXBalance, err error) {
	// make request
	_, _, err = api.doRequest("POST", "/v1/balance", nil, balances)
	return
}
