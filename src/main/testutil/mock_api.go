package testutil

import (
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/stretchr/testify/mock"
)

type MockApi struct {
	mock.Mock
}

func (m *MockApi) FormatAmount(amount float64) string {
	argsCalled := m.Called(amount)
	return argsCalled.String(0)
}

func (m *MockApi) CreateOrder(pair string, fiatAmount float64) (*model.Order, error) {
	argsCalled := m.Called(pair, fiatAmount)
	return argsCalled.Get(0).(*model.Order), argsCalled.Error(1)
}

func (m *MockApi) SubmitOrder(order model.Order) ([]string, error) {
	argsCalled := m.Called(order)
	return argsCalled.Get(0).([]string), argsCalled.Error(1)
}

func (m *MockApi) TransactionStatus(transactionId string) (*krakenapi.Order, error) {
	argsCalled := m.Called(transactionId)
	return argsCalled.Get(0).(*krakenapi.Order), argsCalled.Error(1)
}

func (m *MockApi) IsLive() bool {
	argsCalled := m.Called()
	return argsCalled.Bool(0)
}
