package api

import (
	"fmt"
	"github.com/avast/retry-go"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	"github.com/jackpf/kraken-scheduler/src/main/testutil"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"

	"github.com/stretchr/testify/assert"

	krakenapi "github.com/beldur/kraken-go-api-client"

	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
)

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

type ApiTestSuite struct {
	suite.Suite
	krakenAPI *testutil.MockKrakenApi
	api       Api
}

func (suite *ApiTestSuite) SetupTest() {
	retry.DefaultDelay = 0 * time.Second

	suite.krakenAPI = new(testutil.MockKrakenApi)
	suite.api = NewApi(configmodel.Config{"", "", []configmodel.Schedule{}}, true, true, suite.krakenAPI)
}

func (suite *ApiTestSuite) TestApi_CreateOrder() {
	price := 246.0
	pair := configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}

	suite.krakenAPI.On("Ticker", []string{pair.Name()}).Return(&krakenapi.TickerResponse{
		XXBTZEUR: krakenapi.PairTickerInfo{Close: []string{fmt.Sprintf("%f", price), "0"}},
	}, nil)

	order, err := suite.api.CreateOrder(pair, 123.0)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pair, order.Pair)
	assert.Equal(suite.T(), 123.0, order.FiatAmount)
	assert.Equal(suite.T(), price, order.Price)
	assert.Equal(suite.T(), 0.5, order.Amount())
}

func (suite *ApiTestSuite) TestApi_SubmitOrder() {
	order := model.NewOrder(configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, 123.0, 246.0)
	transactionIds := []string{"1", "2"}
	suite.krakenAPI.On("AddOrder", order.Pair.Name(), "buy", "market", "2", map[string]string{}).Return(
		&krakenapi.AddOrderResponse{TransactionIds: transactionIds},
		nil,
	)

	result, err := suite.api.SubmitOrder(order)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transactionIds, result)
}

func (suite *ApiTestSuite) TestApi_SubmitOrder_NotLive() {
	suite.api = NewApi(configmodel.Config{"", "", []configmodel.Schedule{}}, false, true, suite.krakenAPI)

	order := model.NewOrder(configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, 123.0, 246.0)
	transactionIds := []string{"1", "2"}
	suite.krakenAPI.On("AddOrder", order.Pair.Name(), "buy", "market", "2", map[string]string{"validate": "true"}).Return(
		&krakenapi.AddOrderResponse{TransactionIds: transactionIds},
		nil,
	)

	result, err := suite.api.SubmitOrder(order)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transactionIds, result)
}

func (suite *ApiTestSuite) TestApi_TransactionStatus_Open() {
	transactionId := "test-id"
	order := krakenapi.Order{TransactionID: transactionId}

	suite.krakenAPI.On("OpenOrders", map[string]string{}).Return(&krakenapi.OpenOrdersResponse{
		Count: 1,
		Open:  map[string]krakenapi.Order{transactionId: order},
	},
		nil)

	suite.krakenAPI.On("ClosedOrders", map[string]string{}).Return(&krakenapi.ClosedOrdersResponse{
		Count:  0,
		Closed: map[string]krakenapi.Order{},
	},
		nil)

	result, err := suite.api.TransactionStatus(transactionId)

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *ApiTestSuite) TestApi_TransactionStatus_Closed() {
	transactionId := "test-id"
	order := krakenapi.Order{TransactionID: transactionId}

	suite.krakenAPI.On("OpenOrders", map[string]string{}).Return(&krakenapi.OpenOrdersResponse{
		Count: 0,
		Open:  map[string]krakenapi.Order{},
	},
		nil)

	suite.krakenAPI.On("ClosedOrders", map[string]string{}).Return(&krakenapi.ClosedOrdersResponse{
		Count:  1,
		Closed: map[string]krakenapi.Order{transactionId: order},
	},
		nil)

	result, err := suite.api.TransactionStatus(transactionId)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &order, result)
}

func (suite *ApiTestSuite) TestApi_TransactionStatus_NotFound() {
	transactionId := "test-id"

	suite.krakenAPI.On("OpenOrders", map[string]string{}).Return(&krakenapi.OpenOrdersResponse{
		Count: 0,
		Open:  map[string]krakenapi.Order{},
	},
		nil)

	suite.krakenAPI.On("ClosedOrders", map[string]string{}).Return(&krakenapi.ClosedOrdersResponse{
		Count:  1,
		Closed: map[string]krakenapi.Order{},
	},
		nil)

	result, err := suite.api.TransactionStatus(transactionId)

	assert.EqualError(suite.T(), err, "transaction test-id could not be found in open or closed order history")
	assert.Nil(suite.T(), result)
}

func (suite *ApiTestSuite) TestApiImpl_CheckBalance() {
	suite.krakenAPI.On("Balance").Return(&krakenapi.BalanceResponse{
		ZEUR: 100.0,
		ZUSD: 20.0,
	},
		nil)

	request := []apimodel.BalanceRequest{
		{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Amount: 100.0},
		{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Amount: 200.0},
		{Pair: configmodel.Pair{configmodel.XTZ, configmodel.ZUSD}, Amount: 50.0}}

	response, err := suite.api.CheckBalance(request)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2)
	assert.Contains(suite.T(), response, apimodel.BalanceData{Asset: configmodel.ZEUR, NextPurchaseAmount: 300.0, Balance: 100.0})
	assert.Contains(suite.T(), response, apimodel.BalanceData{Asset: configmodel.ZUSD, NextPurchaseAmount: 50.0, Balance: 20.0})
}
