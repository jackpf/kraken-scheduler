package testutil

import (
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/stretchr/testify/mock"
)

type MockMetrics struct {
	mock.Mock
}

func (m *MockMetrics) LogOrder(pair model.Pair) {
	m.Called(pair)
}

func (m *MockMetrics) LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64, holdingsValue float64) {
	m.Called(pair, amount, fiatAmount, holdings, holdingsValue)
}

func (m *MockMetrics) LogCurrencyBalance(asset model.Asset, holdings float64) {
	m.Called(asset, holdings)
}

func (m *MockMetrics) LogError() {
	m.Called()
}
