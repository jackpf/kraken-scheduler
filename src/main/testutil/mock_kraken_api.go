package testutil

import (
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/stretchr/testify/mock"
)

type MockKrakenApi struct {
	mock.Mock
}

func (m *MockKrakenApi) Ticker(pairs ...string) (*krakenapi.TickerResponse, error) {
	argsCalled := m.Called(pairs)
	return argsCalled.Get(0).(*krakenapi.TickerResponse), argsCalled.Error(1)
}

func (m *MockKrakenApi) AddOrder(pair string, direction string, orderType string, volume string, args map[string]string) (*krakenapi.AddOrderResponse, error) {
	argsCalled := m.Called(pair, direction, orderType, volume, args)
	return argsCalled.Get(0).(*krakenapi.AddOrderResponse), argsCalled.Error(1)
}

func (m *MockKrakenApi) OpenOrders(args map[string]string) (*krakenapi.OpenOrdersResponse, error) {
	argsCalled := m.Called(args)
	return argsCalled.Get(0).(*krakenapi.OpenOrdersResponse), argsCalled.Error(1)
}

func (m *MockKrakenApi) ClosedOrders(args map[string]string) (*krakenapi.ClosedOrdersResponse, error) {
	argsCalled := m.Called(args)
	return argsCalled.Get(0).(*krakenapi.ClosedOrdersResponse), argsCalled.Error(1)
}

func (m *MockKrakenApi) Balance() (*krakenapi.BalanceResponse, error) {
	argsCalled := m.Called()
	return argsCalled.Get(0).(*krakenapi.BalanceResponse), argsCalled.Error(1)
}
