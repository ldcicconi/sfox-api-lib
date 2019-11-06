package sfoxapi

import (
	"github.com/shopspring/decimal"
	"strconv"
)

var (
	SIDE_BUY  = "buy"
	SIDE_SELL = "sell"
)

type OrderStatusResponse struct {
	ID             int64           `json:"id"`
	Quantity       decimal.Decimal `json:"quantity"`
	Price          decimal.Decimal `json:"quantity"`
	Pair           string          `json:"pair"`
	VWAP           decimal.Decimal `json:"vwap"`
	FilledQuantity decimal.Decimal `json:"filled"`
	Status         string          `json:"status"`
}

type NewOrderReqeust struct {
	Quantity      decimal.Decimal `json:"quantity"`
	Price         decimal.Decimal `json:"price"`
	Pair          string          `json:"currency_pair"`
	AlgoID        int             `json:"algorith_id"`
	ClientOrderID string          `json:"client_order_id"`
}

func (api *SFOXAPI) NewOrder(quantity, price decimal.Decimal, algoID int, pair, side string) (orderStatus OrderStatusResponse, err error) {
	// create request body
	reqBody := NewOrderReqeust{
		quantity,
		price,
		pair,
		algoID,
		"",
	}
	// make request
	_, _, err = api.doRequest("POST", "/v1/orders/"+side, reqBody, &orderStatus)
	return
}

func (api *SFOXAPI) OrderStatus(id int) (orderStatus OrderStatusResponse, err error) {
	// make request
	_, _, err = api.doRequest("GET", "/v1/orders/"+strconv.Itoa(id), nil, &orderStatus)
	return
}

func (api *SFOXAPI) GetActiveOrders() (orders []OrderStatusResponse, err error) {
	// make request
	_, _, err = api.doRequest("GET", "/v1/orders/", nil, orders)
	return
}

func (api *SFOXAPI) CancelOrder(id int) (err error) {
	// make request
	_, _, err = api.doRequest("GET", "/v1/orders/", nil, nil)
	return
}
