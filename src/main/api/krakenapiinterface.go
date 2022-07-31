package api

import krakenapi "github.com/beldur/kraken-go-api-client"

type KrakenApiInterface interface {
	Ticker(pairs ...string) (*krakenapi.TickerResponse, error)
	AddOrder(pair string, direction string, orderType string, volume string, args map[string]string) (*krakenapi.AddOrderResponse, error)
	OpenOrders(args map[string]string) (*krakenapi.OpenOrdersResponse, error)
	ClosedOrders(args map[string]string) (*krakenapi.ClosedOrdersResponse, error)
}
