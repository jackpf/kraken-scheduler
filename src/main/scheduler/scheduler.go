package scheduler

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"

	log "github.com/sirupsen/logrus"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/jackpf/kraken-schedule/src/main/config/model"
)

func NewScheduler(appConfig model.Config, api *krakenapi.KrakenAPI, notifier notifier.Notifier) Scheduler {
	return Scheduler{
		config:          appConfig,
		refreshInterval: 1 * time.Minute,
		api:             api,
		cron:            gocron.NewScheduler(time.UTC),
		notifier:        notifier,
	}
}

type Scheduler struct {
	config          model.Config
	refreshInterval time.Duration
	api             *krakenapi.KrakenAPI
	cron            *gocron.Scheduler
	notifier        notifier.Notifier
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

func (s Scheduler) submitOrder(schedule model.Schedule) { // TODO Retry
	currentPrice, err := s.getCurrentPrice(schedule.Pair)
	if err != nil {
		log.Errorf("Unable to fetch price information: %s", err.Error())
		return
	}

	amount := schedule.Amount / *currentPrice

	message := fmt.Sprintf("Buying %+v %s for %+v (%s = %f)", amount, schedule.Pair, schedule.Amount, schedule.Pair, *currentPrice)
	log.Infof(message)
	//err = s.notifier.Send(s.config.NotifyEmailAddress, "kraken-scheduler: Purchase", message)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func (s Scheduler) Run() {
	for {
		for _, schedule := range s.config.Schedules {
			_, err := s.cron.Cron(schedule.Cron).Do(s.submitOrder, schedule)

			if err != nil {
				log.Fatalf("Unable to create cron schedule: %s", err.Error())
			}

			log.Infof("Created cron schedule for %s", schedule.Pair)
		}

		s.cron.StartBlocking()
	}
}
