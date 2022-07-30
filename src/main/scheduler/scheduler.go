package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	"github.com/jackpf/kraken-schedule/src/main/util"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/jackpf/kraken-schedule/src/main/config/model"
)

func NewScheduler(appConfig model.Config, api *krakenapi.KrakenAPI, notifier notifier.Notifier) Scheduler {
	return Scheduler{
		config:             appConfig,
		refreshInterval:    1 * time.Minute,
		api:                api,
		currencyCalculator: NewCurrencyCalculator(api),
		notifier:           notifier,
	}
}

type Scheduler struct {
	config             model.Config
	refreshInterval    time.Duration
	api                *krakenapi.KrakenAPI
	currencyCalculator CurrencyCalculator
	notifier           notifier.Notifier
}

func (s Scheduler) Run() {
	for {
		log.Printf("Checking for outstanding orders for %d %s...", len(s.config.Schedules), util.Pluralise("schedule", len(s.config.Schedules)))

		amount, err := s.currencyCalculator.AmountFor(s.config.Schedules[0].Pair, s.config.Schedules[0].Amount)

		if err != nil {
			log.Fatal(err)
		}

		message := fmt.Sprintf("Buying %+v btc for €%+v", *amount, s.config.Schedules[0].Amount)
		log.Println(message)
		//err = s.notifier.Send(s.config.NotifyEmailAddress, "kraken-scheduler: Purchase", message)
		//
		//if err != nil {
		//	log.Fatal(err)
		//}

		log.Printf("Sleeping for %+v...\n\n", s.refreshInterval)
		time.Sleep(s.refreshInterval)
	}
}
