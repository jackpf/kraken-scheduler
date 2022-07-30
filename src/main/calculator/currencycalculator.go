package calculator

import (
	"fmt"
	"reflect"
	"strconv"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

func NewCurrencyCalculator(api *krakenapi.KrakenAPI) CurrencyCalculator {
	return CurrencyCalculator{api: api}
}

type CurrencyCalculator struct {
	api *krakenapi.KrakenAPI
}

func (cc CurrencyCalculator) getTickerInfo(pair string) (*krakenapi.TickerResponse, error) {
	return cc.api.Ticker(pair)
}

func (cc CurrencyCalculator) calculateLastPrice(tickerInfo krakenapi.PairTickerInfo) (*float32, error) {
	pricePair := tickerInfo.Close

	if len(pricePair) != 2 {
		return nil, fmt.Errorf("expected 2 values, got: %d", len(pricePair))
	}

	price, err := strconv.ParseFloat(pricePair[0], 32)
	if err != nil {
		return nil, err
	}

	price32 := float32(price)

	return &price32, nil
}

func (cc CurrencyCalculator) AmountFor(pair string, fiatAmount float32) (*float32, error) {
	tickerResult, err := cc.getTickerInfo(pair)

	if err != nil {
		return nil, err
	}

	tickerInfo := reflect.ValueOf(*tickerResult).
		FieldByName(pair).
		Interface().(krakenapi.PairTickerInfo)

	price, err := cc.calculateLastPrice(tickerInfo)

	if err != nil {
		return nil, err
	}

	amount := fiatAmount / *price

	return &amount, nil
}
