package sfoxapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type Side int

const (
	SideBuy Side = iota
	SideSell
)

func sideToString(side Side) (sideString string) {
	switch side {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	default:
		return "unknown_side"
	}
}

type OrderStatusResponse struct {
	ID             int64           `json:"id"`
	Quantity       decimal.Decimal `json:"quantity"`
	Price          decimal.Decimal `json:"price"`
	Pair           string          `json:"pair"`
	VWAP           decimal.Decimal `json:"vwap"`
	FilledQuantity decimal.Decimal `json:"filled"`
	Status         string          `json:"status"`
	NetProceeds    decimal.Decimal `json:"net_proceeds"`
}

type NewOrderRequest struct {
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	Pair          string  `json:"currency_pair"`
	AlgoID        int     `json:"algorithm_id"`
	ClientOrderID string  `json:"client_order_id"`
}

func (client *Client) NewOrder(quantity, price decimal.Decimal, algoID int, pair string, side Side) (orderStatus OrderStatusResponse, err error) {
	q, _ := quantity.Float64()
	p, _ := price.Float64()
	if pair[len(pair)-3:] == "usd" && (p*q < 5 || q < 0.001) {
		err = fmt.Errorf("not bigger than the minimum")
		return
	}
	reqBody := NewOrderRequest{
		q,
		p,
		pair,
		algoID,
		"",
	}

	return orderStatus, client.doRequest(http.MethodPost, "/v1/orders/"+sideToString(side), reqBody, &orderStatus, false, true)
}

func (client *Client) OrderStatus(id int64) (orderStatus OrderStatusResponse, err error) { // make request
	return orderStatus, client.doRequest(http.MethodGet, "/v1/orders/"+strconv.FormatInt(id, 10), nil, &orderStatus, false, true)
}

func (client *Client) GetActiveOrders() (orders []OrderStatusResponse, err error) {
	return orders, client.doRequest(http.MethodGet, "/v1/orders/", nil, orders, false, true)
}

func (client *Client) CancelOrder(id int64) (err error) {
	return client.doRequest(http.MethodDelete, "/v1/orders/"+strconv.FormatInt(id, 10), nil, nil, false, true)
}
