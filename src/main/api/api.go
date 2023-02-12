package api

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/jackpf/kraken-scheduler/src/main/util"
	"reflect"
	"strconv"

	krakenapi "github.com/beldur/kraken-go-api-client"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
)

type Api interface {
	CreateOrder(pair configmodel.Pair, fiatAmount float64) (*model.Order, error)
	SubmitOrder(order model.Order) ([]string, error)
	TransactionStatus(transactionId string) (*krakenapi.Order, error)
	CheckBalance(balanceRequests []apimodel.BalanceRequest) ([]apimodel.BalanceData, error)
	IsLive() bool
	IsVerbose() bool
}

func NewApi(appConfig configmodel.Config, live bool, verbose bool, krakenAPI KrakenApiInterface) Api {
	return &ApiImpl{
		config:    appConfig,
		live:      live,
		verbose:   verbose,
		krakenAPI: krakenAPI,
	}
}

type ApiImpl struct {
	config    configmodel.Config
	live      bool
	verbose   bool
	krakenAPI KrakenApiInterface
}

func (a ApiImpl) getCurrentPrice(pair configmodel.Pair) (*float64, error) {
	tickerResult, err := a.krakenAPI.Ticker(pair.Name())

	if err != nil {
		return nil, err
	}

	tickerInfo := reflect.ValueOf(*tickerResult).
		FieldByName(pair.Name()).
		Interface().(krakenapi.PairTickerInfo)

	pricePair := tickerInfo.Close

	if len(pricePair) != 2 {
		return nil, fmt.Errorf("expected 2 values, got: %d", len(pricePair))
	}

	price, err := strconv.ParseFloat(pricePair[0], 32)
	if err != nil {
		return nil, err
	}

	price32 := price

	return &price32, nil
}

func (a ApiImpl) CreateOrder(pair configmodel.Pair, fiatAmount float64) (*model.Order, error) {
	var order model.Order

	if err := retry.Do(func() error {
		currentPrice, err := a.getCurrentPrice(pair)
		if err != nil {
			return err
		}

		order = model.NewOrder(pair, *currentPrice, fiatAmount)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to create order: %s", err.Error())
	}

	return &order, nil
}

func (a ApiImpl) SubmitOrder(order model.Order) ([]string, error) {
	var orderResponse *krakenapi.AddOrderResponse

	if err := retry.Do(func() error {
		data := map[string]string{}
		if !a.live {
			data["validate"] = "true"
		}

		var err error
		if orderResponse, err = a.krakenAPI.AddOrder(order.Pair.Name(), "buy", "market", util.FormatFloat(order.Amount(), 8), data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to submit order: %s", err.Error())
	}

	return orderResponse.TransactionIds, nil
}

func (a ApiImpl) TransactionStatus(transactionId string) (*krakenapi.Order, error) {
	var openOrders *krakenapi.OpenOrdersResponse
	var closedOrders *krakenapi.ClosedOrdersResponse

	if err := retry.Do(func() error {
		var err error
		if openOrders, err = a.krakenAPI.OpenOrders(map[string]string{}); err != nil {
			return err
		}
		if closedOrders, err = a.krakenAPI.ClosedOrders(map[string]string{}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to check transaction status: %s", err.Error())
	}

	if _, isOpen := openOrders.Open[transactionId]; isOpen {
		return nil, nil
	} else if order, isClosed := closedOrders.Closed[transactionId]; isClosed {
		return &order, nil
	} else {
		return nil, fmt.Errorf("transaction %s could not be found in open or closed order history", transactionId)
	}
}

func (a ApiImpl) CheckBalance(balanceRequests []apimodel.BalanceRequest) ([]apimodel.BalanceData, error) {
	var balance *krakenapi.BalanceResponse

	if err := retry.Do(func() error {
		var err error
		if balance, err = a.krakenAPI.Balance(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to check balance: %s", err.Error())
	}

	totalToPurchase := make(map[configmodel.Asset]float64)

	for _, balanceRequest := range balanceRequests {
		totalToPurchase[balanceRequest.Pair.Second] += balanceRequest.Amount
	}

	var balanceData []apimodel.BalanceData

	for asset, amount := range totalToPurchase {
		balanceInCurrency := reflect.ValueOf(*balance).
			FieldByName(asset.NormalisedName).
			Interface().(float64)

		balanceData = append(balanceData, apimodel.BalanceData{
			Asset:              asset,
			NextPurchaseAmount: amount,
			Balance:            balanceInCurrency,
		})
	}

	return balanceData, nil
}

func (a ApiImpl) IsLive() bool {
	return a.live
}

func (a ApiImpl) IsVerbose() bool {
	return a.verbose
}
