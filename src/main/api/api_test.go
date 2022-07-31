package api

import (
	"fmt"
	"testing"

	"github.com/jackpf/kraken-schedule/src/main/scheduler/model"

	"github.com/stretchr/testify/assert"

	krakenapi "github.com/beldur/kraken-go-api-client"

	configmodel "github.com/jackpf/kraken-schedule/src/main/config/model"
	"github.com/stretchr/testify/mock"
)

type MockKrakenApi struct {
	mock.Mock
}

func (k *MockKrakenApi) Ticker(pairs ...string) (*krakenapi.TickerResponse, error) {
	argsCalled := k.Called(pairs)
	return argsCalled.Get(0).(*krakenapi.TickerResponse), argsCalled.Error(1)
}

func (k *MockKrakenApi) AddOrder(pair string, direction string, orderType string, volume string, args map[string]string) (*krakenapi.AddOrderResponse, error) {
	argsCalled := k.Called(pair, direction, orderType, volume, args)
	return argsCalled.Get(0).(*krakenapi.AddOrderResponse), argsCalled.Error(1)
}

func TestApi_FormatAmount(t *testing.T) {
	krakenAPI := new(MockKrakenApi)
	api := NewApi(configmodel.Config{"", []configmodel.Schedule{}}, true, krakenAPI)

	result := api.FormatAmount(12.3456789)

	assert.Equal(t, "12.3457", result)
}

func TestApi_CreateOrder(t *testing.T) {
	krakenAPI := new(MockKrakenApi)
	api := NewApi(configmodel.Config{"", []configmodel.Schedule{}}, true, krakenAPI)

	price := float32(246.0)
	pair := "XXBTZEUR" // Must be a real pair due to reflection use

	krakenAPI.On("Ticker", []string{pair}).Return(&krakenapi.TickerResponse{
		XXBTZEUR: krakenapi.PairTickerInfo{Close: []string{fmt.Sprintf("%f", price), "0"}},
	}, nil)

	order, err := api.CreateOrder(configmodel.Schedule{"", pair, 123.0})

	assert.NoError(t, err)
	assert.Equal(t, pair, order.Pair)
	assert.Equal(t, float32(123.0), order.FiatAmount)
	assert.Equal(t, price, order.Price)
	assert.Equal(t, float32(0.5), order.Amount())
}

func TestApi_ValidateOrder(t *testing.T) {
	krakenAPI := new(MockKrakenApi)
	api := NewApi(configmodel.Config{"", []configmodel.Schedule{}}, true, krakenAPI)

	resultValid := api.ValidateOrder(model.NewOrder("test-pair", 123.0, 456.0))
	assert.NoError(t, resultValid)

	resultInvalid := api.ValidateOrder(model.NewOrder("test-pair", 123.0, 0.0))
	assert.EqualError(t, resultInvalid, "order amount too small: 0.000000")
}

func TestApi_SubmitOrder(t *testing.T) {
	krakenAPI := new(MockKrakenApi)
	api := NewApi(configmodel.Config{"", []configmodel.Schedule{}}, true, krakenAPI)

	order := model.NewOrder("test-pair", 123.0, 246.0)
	transactionIds := []string{"1", "2"}
	krakenAPI.On("AddOrder", order.Pair, "buy", "market", "2.0000", map[string]string{}).Return(
		&krakenapi.AddOrderResponse{TransactionIds: transactionIds},
		nil,
	)

	result, err := api.SubmitOrder(order)

	assert.NoError(t, err)
	assert.Equal(t, transactionIds, result)
}

func TestApi_SubmitOrder_NotLive(t *testing.T) {
	krakenAPI := new(MockKrakenApi)
	api := NewApi(configmodel.Config{"", []configmodel.Schedule{}}, false, krakenAPI)

	order := model.NewOrder("test-pair", 123.0, 246.0)
	transactionIds := []string{"1", "2"}
	krakenAPI.On("AddOrder", order.Pair, "buy", "market", "2.0000", map[string]string{"validate": "true"}).Return(
		&krakenapi.AddOrderResponse{TransactionIds: transactionIds},
		nil,
	)

	result, err := api.SubmitOrder(order)

	assert.NoError(t, err)
	assert.Equal(t, transactionIds, result)
}