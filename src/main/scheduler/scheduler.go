package scheduler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackpf/kraken-schedule/src/main/scheduler/model"

	"github.com/go-co-op/gocron"

	log "github.com/sirupsen/logrus"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	krakenapi "github.com/beldur/kraken-go-api-client"
	configmodel "github.com/jackpf/kraken-schedule/src/main/config/model"
)

func NewScheduler(appConfig configmodel.Config, live bool, api *krakenapi.KrakenAPI, notifier *notifier.Notifier) Scheduler {
	return Scheduler{
		config:          appConfig,
		live:            live,
		refreshInterval: 1 * time.Minute,
		api:             api,
		cron:            gocron.NewScheduler(time.UTC),
		notifier:        notifier,
	}
}

type Scheduler struct {
	config          configmodel.Config
	live            bool
	refreshInterval time.Duration
	api             *krakenapi.KrakenAPI
	cron            *gocron.Scheduler
	notifier        *notifier.Notifier
}

func (s Scheduler) liveLogTag() string {
	if s.live {
		return "LIVE"
	}
	return "TEST"
}

func (s Scheduler) getCurrentPrice(pair string) (*float32, error) {
	tickerResult, err := s.api.Ticker(pair)

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

func (s Scheduler) formatAmount(amount float32) string {
	return fmt.Sprintf("%.4f", amount)
}

func (s Scheduler) createOrder(schedule configmodel.Schedule) (*model.Order, error) { // TODO Retry
	currentPrice, err := s.getCurrentPrice(schedule.Pair)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch price information: %s", err.Error())
	}

	order := model.NewOrder(schedule.Pair, *currentPrice, schedule.Amount)

	return &order, nil
}
func (s Scheduler) validateOrder(order model.Order) error {
	if order.Amount() < 0.0001 {
		return fmt.Errorf("order amount too small: %f", order.Amount())
	}

	return nil
}

// TODO Check order status & send confirmation
func (s Scheduler) submitOrder(order model.Order) error { // TODO Retry
	log.Infof("[%s] Ordering %s %s for %+v (%s = %f)...", s.liveLogTag(), s.formatAmount(order.Amount()), order.Pair, order.FiatAmount, order.Pair, order.Price)

	data := map[string]string{}
	if !s.live {
		data["validate"] = "true"
	}

	orderResponse, err := s.api.AddOrder(order.Pair, "buy", "market", s.formatAmount(order.Amount()), data)

	if err != nil {
		return err
	}

	transactionIdsString := strings.Join(orderResponse.TransactionIds[:], ", ")
	if !s.live {
		transactionIdsString = "<no transaction IDs for test orders>"
	}

	log.Infof("[%s] Order placed: %s", s.liveLogTag(), transactionIdsString)

	return nil
}

func (s Scheduler) notifyOrder(order model.Order) error {
	if s.notifier == nil {
		log.Warn("Notifications not configured, not notifying")
		return nil
	}

	message := fmt.Sprintf("[%s] Ordered %s %s for %+v (%s = %f)...", s.liveLogTag(), s.formatAmount(order.Amount()), order.Pair, order.FiatAmount, order.Pair, order.Price)
	return (*s.notifier).Send(s.config.NotifyEmailAddress, "kraken-scheduler: Purchase", message)
}

func (s Scheduler) process(schedule configmodel.Schedule) {
	order, err := s.createOrder(schedule)
	if err != nil {
		log.Errorf("Unable to create order: %s", err.Error())
		return
	}

	err = s.validateOrder(*order)
	if err != nil {
		log.Errorf("Unable to validate order: %s", err.Error())
		return
	}

	err = s.submitOrder(*order)
	if err != nil {
		log.Errorf("Unable to submit order: %s", err.Error())
		return
	}

	err = s.notifyOrder(*order)
	if err != nil {
		log.Errorf("Unable to notify order: %s", err.Error())
		return
	}
}

func (s Scheduler) Run() {
	for {
		for _, schedule := range s.config.Schedules {
			_, err := s.cron.Cron(schedule.Cron).Do(s.process, schedule)

			if err != nil {
				log.Fatalf("Unable to create cron schedule: %s", err.Error())
			}

			log.Infof("Created cron schedule for %s", schedule.Pair)
		}

		s.cron.StartBlocking()
	}
}
