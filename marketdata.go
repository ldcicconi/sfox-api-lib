package sfoxapi

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/valyala/fastjson"
)

var One = decimal.NewFromFloat(1)
var Two = decimal.NewFromFloat(2)

type SfoxOrderbook struct {
	Bids []Offer `json:"bids"`
	Asks []Offer `json:"asks"`
}

func (ob *SfoxOrderbook) MidPrice() (decimal.Decimal, error) {
	if len(ob.Bids) < 1 || len(ob.Asks) < 1 {
		return decimal.Zero, errors.New("bids or asks not long enough")
	}

	return ob.Asks[0].Price.Add(ob.Bids[0].Price).Div(Two), nil
}

func (ob *SfoxOrderbook) WeightedMidPriceSimple() (decimal.Decimal, error) {
	if len(ob.Bids) < 1 || len(ob.Asks) < 1 {
		return decimal.Zero, errors.New("bids or asks not long enough")
	}

	sumQuantities := ob.Bids[0].Quantity.Add(ob.Asks[0].Quantity)
	if sumQuantities.IsZero() {
		return decimal.Zero, errors.New("zero quantity ???")
	}

	imbalance := ob.Bids[0].Quantity.Div(sumQuantities)

	firstTerm := imbalance.Mul(ob.Asks[0].Price)
	secondTerm := One.Sub(imbalance).Mul(ob.Bids[0].Quantity)

	return firstTerm.Add(secondTerm), nil
}

type Offer struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
	Exchange string
}

func (api *SFOXAPI) GetOrderbook(pair string) (ob SfoxOrderbook, err error) {
	bodyBytes, _, err := api.doRequest("GET", "/v1/markets/orderbook/"+pair, nil, nil, false, false)
	if err != nil {
		return ob, err
	}

	orderbook, err := NewSFOXOrderbookFromJSON(bodyBytes)

	return *orderbook, err
}

func NewSFOXOrderbookFromJSON(jsonBytes []byte) (o *SfoxOrderbook, err error) {
	var p fastjson.Parser
	bookJSON, err := p.ParseBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	// fill bids
	bids := bookJSON.GetArray("bids")
	asks := bookJSON.GetArray("asks")

	return &SfoxOrderbook{
		Bids: FillOffersFromJSON(bids),
		Asks: FillOffersFromJSON(asks),
	}, nil
}

func FillOffersFromJSON(offers []*fastjson.Value) (ret []Offer) {
	for _, offer := range offers {
		offerArray := offer.GetArray()
		exBytes, _ := offerArray[2].StringBytes()
		o := Offer{
			Price:    decimal.NewFromFloat(offerArray[0].GetFloat64()),
			Quantity: decimal.NewFromFloat(offerArray[1].GetFloat64()),
			Exchange: string(exBytes),
		}
		ret = append(ret, o)
	}
	return
}
