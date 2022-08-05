package api

import (
	"fmt"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	"github.com/jackpf/kraken-scheduler/src/main/testutil"
	"testing"

	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"

	"github.com/stretchr/testify/assert"

	krakenapi "github.com/beldur/kraken-go-api-client"

	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
)

func TestApi_CreateOrder(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, true, krakenAPI)

	price := 246.0
	pair := configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}

	krakenAPI.On("Ticker", []string{pair.Name()}).Return(&krakenapi.TickerResponse{
		XXBTZEUR: krakenapi.PairTickerInfo{Close: []string{fmt.Sprintf("%f", price), "0"}},
	}, nil)

	order, err := api.CreateOrder(pair, 123.0)

	assert.NoError(t, err)
	assert.Equal(t, pair, order.Pair)
	assert.Equal(t, 123.0, order.FiatAmount)
	assert.Equal(t, price, order.Price)
	assert.Equal(t, 0.5, order.Amount())
}

func TestApi_SubmitOrder(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, true, krakenAPI)

	order := model.NewOrder(configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, 123.0, 246.0)
	transactionIds := []string{"1", "2"}
	krakenAPI.On("AddOrder", order.Pair.Name(), "buy", "market", "2.00000000", map[string]string{}).Return(
		&krakenapi.AddOrderResponse{TransactionIds: transactionIds},
		nil,
	)

	result, err := api.SubmitOrder(order)

	assert.NoError(t, err)
	assert.Equal(t, transactionIds, result)
}

func TestApi_SubmitOrder_NotLive(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, false, krakenAPI)

	order := model.NewOrder(configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, 123.0, 246.0)
	transactionIds := []string{"1", "2"}
	krakenAPI.On("AddOrder", order.Pair.Name(), "buy", "market", "2.00000000", map[string]string{"validate": "true"}).Return(
		&krakenapi.AddOrderResponse{TransactionIds: transactionIds},
		nil,
	)

	result, err := api.SubmitOrder(order)

	assert.NoError(t, err)
	assert.Equal(t, transactionIds, result)
}

func TestApi_TransactionStatus_Open(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, false, krakenAPI)

	transactionId := "test-id"
	order := krakenapi.Order{TransactionID: transactionId}

	krakenAPI.On("OpenOrders", map[string]string{}).Return(&krakenapi.OpenOrdersResponse{
		Count: 1,
		Open:  map[string]krakenapi.Order{transactionId: order},
	},
		nil)

	krakenAPI.On("ClosedOrders", map[string]string{}).Return(&krakenapi.ClosedOrdersResponse{
		Count:  0,
		Closed: map[string]krakenapi.Order{},
	},
		nil)

	result, err := api.TransactionStatus(transactionId)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestApi_TransactionStatus_Closed(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, false, krakenAPI)

	transactionId := "test-id"
	order := krakenapi.Order{TransactionID: transactionId}

	krakenAPI.On("OpenOrders", map[string]string{}).Return(&krakenapi.OpenOrdersResponse{
		Count: 0,
		Open:  map[string]krakenapi.Order{},
	},
		nil)

	krakenAPI.On("ClosedOrders", map[string]string{}).Return(&krakenapi.ClosedOrdersResponse{
		Count:  1,
		Closed: map[string]krakenapi.Order{transactionId: order},
	},
		nil)

	result, err := api.TransactionStatus(transactionId)

	assert.NoError(t, err)
	assert.Equal(t, &order, result)
}

func TestApi_TransactionStatus_NotFound(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, false, krakenAPI)

	transactionId := "test-id"

	krakenAPI.On("OpenOrders", map[string]string{}).Return(&krakenapi.OpenOrdersResponse{
		Count: 0,
		Open:  map[string]krakenapi.Order{},
	},
		nil)

	krakenAPI.On("ClosedOrders", map[string]string{}).Return(&krakenapi.ClosedOrdersResponse{
		Count:  1,
		Closed: map[string]krakenapi.Order{},
	},
		nil)

	result, err := api.TransactionStatus(transactionId)

	assert.EqualError(t, err, "transaction test-id could not be found in open or closed order history")
	assert.Nil(t, result)
}

func TestApiImpl_CheckBalance(t *testing.T) {
	krakenAPI := new(testutil.MockKrakenApi)
	api := NewApi(configmodel.Config{[]configmodel.Schedule{}}, false, krakenAPI)

	krakenAPI.On("Balance").Return(&krakenapi.BalanceResponse{
		ZEUR: 100.0,
		ZUSD: 20.0,
	},
		nil)

	request := []apimodel.BalanceRequest{
		{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Amount: 100.0},
		{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Amount: 200.0},
		{Pair: configmodel.Pair{configmodel.XTZ, configmodel.ZUSD}, Amount: 50.0}}

	response, err := api.CheckBalance(request)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Contains(t, response, apimodel.BalanceData{Asset: configmodel.ZEUR, NextPurchaseAmount: 300.0, Balance: 100.0})
	assert.Contains(t, response, apimodel.BalanceData{Asset: configmodel.ZUSD, NextPurchaseAmount: 50.0, Balance: 20.0})
}
