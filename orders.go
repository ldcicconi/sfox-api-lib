package sfoxapi

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

var (
	SIDE_BUY  = "buy"
	SIDE_SELL = "sell"
)

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

type NewOrderReqeust struct {
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	Pair          string  `json:"currency_pair"`
	AlgoID        int     `json:"algorithm_id"`
	ClientOrderID string  `json:"client_order_id"`
}

func (api *SFOXAPI) NewOrder(quantity, price decimal.Decimal, algoID int, pair, side string) (orderStatus OrderStatusResponse, err error) {
	// create request body
	q, _ := quantity.Float64()
	p, _ := price.Float64()
	if pair[len(pair)-3:] == "usd" && (p*q < 5 || q < 0.001) {
		err = fmt.Errorf("not bigger than the minimum")
		return
	}
	reqBody := NewOrderReqeust{
		q,
		p,
		pair,
		algoID,
		"",
	}
	// make request
	_, _, err = api.doRequest("POST", "/v1/orders/"+side, reqBody, &orderStatus, false, true)
	api.reportError(CreateOrderKey, err)
	return
}

func (api *SFOXAPI) OrderStatus(id int64) (orderStatus OrderStatusResponse, err error) {
	// make request
	_, _, err = api.doRequest("GET", "/v1/orders/"+strconv.FormatInt(id, 10), nil, &orderStatus, false, true)
	api.reportError(OrderStatusKey, err)
	return
}

func (api *SFOXAPI) GetActiveOrders() (orders []OrderStatusResponse, err error) {
	// make request
	_, _, err = api.doRequest("GET", "/v1/orders/", nil, orders, false, true)
	api.reportError(GetOpenOrdersKey, err)
	return
}

func (api *SFOXAPI) CancelOrder(id int64) (err error) {
	// make request
	_, _, err = api.doRequest("DELETE", "/v1/orders/"+strconv.FormatInt(id, 10), nil, nil, false, true)
	api.reportError(CancelOrderKey, err)
	return
}
