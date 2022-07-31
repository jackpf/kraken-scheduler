package api

import (
	"fmt"
	"reflect"
	"strconv"

	krakenapi "github.com/beldur/kraken-go-api-client"
	configmodel "github.com/jackpf/kraken-schedule/src/main/config/model"
	"github.com/jackpf/kraken-schedule/src/main/scheduler/model"
)

func NewApi(appConfig configmodel.Config, live bool, krakenAPI KrakenApiInterface) *Api {
	return &Api{
		config:    appConfig,
		Live:      live,
		krakenAPI: krakenAPI,
	}
}

type Api struct {
	config    configmodel.Config
	Live      bool
	krakenAPI KrakenApiInterface
}

func (a Api) getCurrentPrice(pair string) (*float32, error) {
	tickerResult, err := a.krakenAPI.Ticker(pair)

	if err != nil {
		return nil, err
	}

	tickerInfo := reflect.ValueOf(*tickerResult).
		FieldByName(pair).
		Interface().(krakenapi.PairTickerInfo)

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

func (a Api) FormatAmount(amount float32) string {
	return fmt.Sprintf("%.4f", amount)
}

func (a Api) CreateOrder(schedule configmodel.Schedule) (*model.Order, error) { // TODO Retry
	currentPrice, err := a.getCurrentPrice(schedule.Pair)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch price information: %s", err.Error())
	}

	order := model.NewOrder(schedule.Pair, *currentPrice, schedule.Amount)

	return &order, nil
}

func (a Api) ValidateOrder(order model.Order) error {
	if order.Amount() < 0.0001 {
		return fmt.Errorf("order amount too small: %f", order.Amount())
	}

	return nil
}

// TODO Check order status & send confirmation
func (a Api) SubmitOrder(order model.Order) ([]string, error) { // TODO Retry
	data := map[string]string{}
	if !a.Live {
		data["validate"] = "true"
	}

	orderResponse, err := a.krakenAPI.AddOrder(order.Pair, "buy", "market", a.FormatAmount(order.Amount()), data)

	if err != nil {
		return nil, err
	}

	return orderResponse.TransactionIds, nil
}
