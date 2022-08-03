package scheduler

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackpf/kraken-scheduler/src/main/util"

	"github.com/jackpf/kraken-scheduler/src/main/notificationtemplates"

	"github.com/jackpf/kraken-scheduler/src/main/api"

	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"

	"github.com/go-co-op/gocron"

	log "github.com/sirupsen/logrus"

	"github.com/jackpf/kraken-scheduler/src/main/notifier"

	krakenapi "github.com/beldur/kraken-go-api-client"
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
)

func NewScheduler(appConfig configmodel.Config, api api.Api, notifiers []*notifier.Notifier) Scheduler {
	return Scheduler{
		config:    appConfig,
		api:       api,
		cron:      gocron.NewScheduler(time.Now().Location()),
		notifiers: notifiers,
	}
}

type Scheduler struct {
	config    configmodel.Config
	api       api.Api
	cron      *gocron.Scheduler
	notifiers []*notifier.Notifier
	jobs      []struct {
		configmodel.Schedule
		*gocron.Job
	}
	// State & mutex required for printing console output/loading bars correctly
	startTime time.Time
	mutex     sync.Mutex
	jobRuns   uint64
}

func (s *Scheduler) liveLogTag() string {
	if s.api.IsLive() {
		return "LIVE"
	}
	return "TEST"
}

func (s *Scheduler) notifyOrder(order model.Order, transactionIds []string) []error {
	if len(s.notifiers) == 0 {
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

	var errors []error
	for _, notifier := range s.notifiers {
		var err = (*notifier).Send(notification.Subject(), notification.Body())
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *Scheduler) notifyCompletedTrade(order model.Order, completedOrder krakenapi.Order, transactionId string) []error {
	if len(s.notifiers) == 0 {
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

	var errors []error
	for _, notifier := range s.notifiers {
		var err = (*notifier).Send(notification.Subject(), notification.Body())
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *Scheduler) process(schedule configmodel.Schedule) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	atomic.AddUint64(&s.jobRuns, 1)

	order, err := s.api.CreateOrder(schedule)
	if err != nil {
		log.Errorf("Unable to create order: %s", err.Error())
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

	var errors = s.notifyOrder(*order, transactionIds)

	if errors != nil {
		for _, err := range errors {
			log.Errorf("Unable to notify of completed order: %s", err.Error())
		}
	}

	for _, transactionId := range transactionIds {
		for { // TODO perform in background & have max attempts
			completedOrder, err := s.api.TransactionStatus(transactionId)

			if err != nil {
				log.Errorf("Unable to check transaction status: %s", err.Error())
			}

			if completedOrder != nil {
				log.Infof("Order %s was successfully completed", transactionId)

				var errors = s.notifyCompletedTrade(*order, *completedOrder, transactionId)

				if errors != nil {
					for _, err := range errors {
						log.Errorf("Unable to notify of completed order: %s", err.Error())
					}
				}
				break
			} else {
				log.Infof("Order %s is pending...", transactionId)
				time.Sleep(1 * time.Second)
			}
		}
	}

	job := s.findJob(schedule)
	if job != nil {
		log.Infof("Next purchase for %s will occur at %+v", job.Pair, job.NextRun())
	}
}

func (s *Scheduler) validateSchedule(schedule configmodel.Schedule) error {
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

func (s *Scheduler) findJob(schedule configmodel.Schedule) *struct {
	configmodel.Schedule
	*gocron.Job
} {
	for _, job := range s.jobs {
		if job.Schedule == schedule {
			return &job
		}
	}

	return nil
}

func (s *Scheduler) printJobTimeDiffs() {
	first := true
	lastJobRuns := uint64(0)

	for {
		s.mutex.Lock()

		if !first && lastJobRuns == s.jobRuns {
			util.ClearConsoleLines(len(s.jobs))
		}
		first = false

		for _, job := range s.jobs {
			lastRunTime := job.LastRun().Unix()
			if job.RunCount() == 0 {
				lastRunTime = s.startTime.Unix()
			}

			completedRatio := float64(time.Now().Unix()-lastRunTime) / float64(job.NextRun().Unix()-lastRunTime)

			logOutput := util.PadLine(fmt.Sprintf("Purchasing %s in %s", job.Pair, util.PrettyDuration(time.Until(job.NextRun()))), 60)
			fmt.Printf("%s%s\n", logOutput, util.ProgressBar(completedRatio, 30))
		}
		lastJobRuns = s.jobRuns

		s.mutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (s *Scheduler) Run() {
	s.startTime = time.Now()

	for _, schedule := range s.config.Schedules {
		err := s.validateSchedule(schedule)
		if err != nil {
			log.Fatalf("Invalid schedule: %s", err.Error())
		}

		job, err := s.cron.Cron(schedule.Cron).Do(s.process, schedule)
		if err != nil {
			log.Fatalf("Unable to create cron schedule: %s", err.Error())
		}

		s.jobs = append(s.jobs, struct {
			configmodel.Schedule
			*gocron.Job
		}{schedule, job})
	}

	// Jobs don't have next run information until the scheduler is started
	// Start async, then block after
	s.cron.StartAsync()

	for _, job := range s.jobs {
		log.Infof("Created schedule for %s, purchase will occur at %+v", job.Pair, job.NextRun())
	}

	go s.printJobTimeDiffs()

	s.cron.StartBlocking()
}
