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

type ExecutionAlgorithm int

const (
	ExecutionAlgorithmMarket ExecutionAlgorithm = iota
	ExecutionAlgorithmInstant
	ExecutionAlgorithmSimple
	ExecutionAlgorithmSmart
	ExecutionAlgorithmLimit
	ExecutionAlgorithmGorilla
	ExecutionAlgorithmTortoise
	ExecutionAlgorithmHare
	ExecutionAlgorithmStop
	ExecutionAlgorithmPolarBear
	ExecutionAlgorithmSniper
	ExecutionAlgorithmTWAP
	ExecutionAlgorithmTrailingStop
	ExecutionAlgorithmImmediateOrCancel
)

func algorithmIDFromExecutionAlgorithm(algorithm ExecutionAlgorithm) (algorithmID int) {
	switch algorithm {
	case ExecutionAlgorithmMarket:
		return 100
	case ExecutionAlgorithmSmart:
		return 200
	case ExecutionAlgorithmTWAP:
		return 307
	default:
		return 0
	}
}

type NewOrderRequest struct {
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	Pair          string  `json:"currency_pair"`
	AlgoID        int     `json:"algorithm_id"`
	ClientOrderID string  `json:"client_order_id"`
}

func (client *Client) NewOrder(quantity, price decimal.Decimal, pair string, side Side, algorithm ExecutionAlgorithm, clientOrderID string) (orderStatus OrderStatusResponse, err error) {
	var (
		q, _    = quantity.Float64()
		p, _    = price.Float64()
		reqBody = NewOrderRequest{
			q,
			p,
			pair,
			algorithmIDFromExecutionAlgorithm(algorithm),
			clientOrderID,
		}
	)

	if pair[len(pair)-3:] == "usd" && (p*q < 5 || q < 0.001) {
		err = fmt.Errorf("not bigger than the minimum")
		return
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
