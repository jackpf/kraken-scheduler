package scheduler

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackpf/kraken-schedule/src/main/notificationtemplates"

	"github.com/jackpf/kraken-schedule/src/main/api"

	"github.com/jackpf/kraken-schedule/src/main/scheduler/model"

	"github.com/go-co-op/gocron"

	log "github.com/sirupsen/logrus"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	krakenapi "github.com/beldur/kraken-go-api-client"
	configmodel "github.com/jackpf/kraken-schedule/src/main/config/model"
)

func NewScheduler(appConfig configmodel.Config, api api.Api, notifier *notifier.Notifier) Scheduler {
	return Scheduler{
		config:          appConfig,
		api:             api,
		refreshInterval: 1 * time.Minute,
		cron:            gocron.NewScheduler(time.UTC),
		notifier:        notifier,
	}
}

type Scheduler struct {
	config          configmodel.Config
	api             api.Api
	refreshInterval time.Duration
	cron            *gocron.Scheduler
	notifier        *notifier.Notifier
}

func (s Scheduler) liveLogTag() string {
	if s.api.IsLive() {
		return "LIVE"
	}
	return "TEST"
}

func (s Scheduler) notifyOrder(order model.Order, transactionIds []string) error {
	if s.notifier == nil || s.config.NotifyEmailAddress == "" {
		log.Warn("Notifications not configured, not notifying")
		return nil
	}

	notification := notificationtemplates.NewOrderNotification(
		s.api.IsLive(),
		order.Pair,
		order.Amount(),
		order.FiatAmount,
		order.Price,
		transactionIds,
	)

	return (*s.notifier).Send(s.config.NotifyEmailAddress, notification.Subject(), notification.Body())
}

func (s Scheduler) notifyCompletedTrade(order model.Order, completedOrder krakenapi.Order, transactionId string) error {
	if s.notifier == nil || s.config.NotifyEmailAddress == "" {
		log.Warn("Notifications not configured, not notifying")
		return nil
	}

	notification := notificationtemplates.NewPurchaseNotification(
		order.Pair,
		order.Amount(),
		order.FiatAmount,
		transactionId,
		completedOrder,
	)

	return (*s.notifier).Send(s.config.NotifyEmailAddress, notification.Subject(), notification.Body())
}

func (s Scheduler) process(schedule configmodel.Schedule) {
	order, err := s.api.CreateOrder(schedule)
	if err != nil {
		log.Errorf("Unable to create order: %s", err.Error())
		return
	}

	err = s.api.ValidateOrder(*order)
	if err != nil {
		log.Errorf("Unable to validate order: %s", err.Error())
		return
	}

	log.Infof("[%s] Ordering %s %s for %+v (%s = %f)...", s.liveLogTag(), s.api.FormatAmount(order.Amount()), order.Pair, order.FiatAmount, order.Pair, order.Price)
	transactionIds, err := s.api.SubmitOrder(*order)
	if err != nil {
		log.Errorf("Unable to submit order: %s", err.Error())
		return
	}

	transactionIdsString := strings.Join(transactionIds[:], ", ")
	if !s.api.IsLive() {
		transactionIdsString = "<no transaction IDs for test orders>"
	}

	log.Infof("[%s] Order placed: %s", s.liveLogTag(), transactionIdsString)

	err = s.notifyOrder(*order, transactionIds)
	if err != nil {
		log.Errorf("Unable to notify of order: %s", err.Error())
	}

	for _, transactionId := range transactionIds {
		completedOrder, err := s.api.TransactionStatus(transactionId) // TODO needs to be retried & performed in background

		if err != nil {
			log.Errorf("Unable to check transaction status: %s", err.Error())
		}

		if completedOrder != nil {
			err = s.notifyCompletedTrade(*order, *completedOrder, transactionId)
			if err != nil {
				log.Errorf("Unable to notify of completed order: %s", err.Error())
			}
		}
	}
}

func (s Scheduler) validateSchedule(schedule configmodel.Schedule) error {
	// Ensure pair is valid
	if !reflect.ValueOf(krakenapi.AssetPairsResponse{}).
		FieldByName(schedule.Pair).IsValid() {
		return fmt.Errorf("%s is not a valid asset pair", schedule.Pair)
	}

	// Ensure valid amount
	if schedule.Amount <= 0.0 {
		return fmt.Errorf("purchase amount must be >= 0, got %f", schedule.Amount)
	}

	return nil
}

func (s Scheduler) Run() {
	for {
		for _, schedule := range s.config.Schedules {
			err := s.validateSchedule(schedule)
			if err != nil {
				log.Fatalf("Invalid schedule: %s", err.Error())
			}

			_, err = s.cron.Cron(schedule.Cron).Do(s.process, schedule)
			if err != nil {
				log.Fatalf("Unable to create cron schedule: %s", err.Error())
			}

			log.Infof("Created cron schedule for %s", schedule.Pair)
		}

		s.cron.StartBlocking()
	}
}
