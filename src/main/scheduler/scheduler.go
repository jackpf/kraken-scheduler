package scheduler

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/tasks"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackpf/kraken-scheduler/src/main/util"

	"github.com/jackpf/kraken-scheduler/src/main/notifications"

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
		config: appConfig,
		api:    api,
		cron:   gocron.NewScheduler(time.Now().Location()),
		tasks: []tasks.Task{
			tasks.NewCreateOrderTask(api),
			tasks.NewSubmitOrderTask(api),
			tasks.NewCheckOrderStatusTask(api),
		},
		notifiers: notifiers,
	}
}

type Scheduler struct {
	config    configmodel.Config
	api       api.Api
	cron      *gocron.Scheduler
	tasks     []tasks.Task
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

func (s *Scheduler) logErrors(errs []error) {
	if errs != nil {
		for _, err := range errs {
			log.Error(err.Error())
		}
	}
}

func (s *Scheduler) notifyError(taskData model.TaskData, err error) []error {
	if len(s.notifiers) == 0 {
		log.Warn("Notifications not configured, not notifying")
		return nil
	}

	notification := notifications.NewErrorNotification(
		taskData.Schedule,
		err,
	)

	return s.notify(notification)
}

func (s *Scheduler) notify(notification notifications.Notification) []error {
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

	taskData := model.TaskData{Schedule: schedule}

	for _, task := range s.tasks {
		taskData, err := task.Run(&taskData)
		if err != nil {
			s.notifyError(*taskData, err)
		}

		notifications, errs := task.Notifications(taskData)
		for _, err := range errs {
			s.logErrors(s.notifyError(*taskData, err))
		}
		for _, notification := range notifications {
			s.logErrors(s.notify(notification))
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

func (s *Scheduler) runUi() {
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

	go s.runUi()

	s.cron.StartBlocking()
}
