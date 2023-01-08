package sfoxapi

import (
	"github.com/shopspring/decimal"
	"net/http"
)

type SFOXBalance struct {
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
}

func (client *Client) GetBalances() (balances []SFOXBalance, err error) {
	return balances, client.doRequest(http.MethodGet, "/v1/users/balance", nil, &balances, false, true)
}
