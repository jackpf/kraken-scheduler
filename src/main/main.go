package main

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/jackpf/kraken-scheduler/src/main/api"
	"github.com/jackpf/kraken-scheduler/src/main/config"
	"github.com/jackpf/kraken-scheduler/src/main/metrics"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/jackpf/kraken-scheduler/src/main/notifier"

	"github.com/alexflint/go-arg"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

func main() {
	var args struct {
		ConfigFile              string `arg:"--config,required" help:"Scheduler configuration file"`
		EmailCredentialsFile    string `arg:"--email-credentials" help:"Your google OAuth email-credentials.json file (optional)"`
		TelegramCredentialsFile string `arg:"--telegram-credentials" help:"Your telegram ChatID and Token telegram-credentials.json file (optional)"`
		IsVerbose               bool   `arg:"--verbose" help:"Sends more detailed notifications (order submissions, order logs etc.)"`
		IsLive                  bool   `arg:"--live" default:"false" help:"Set to true to execute real orders"`
	}
	arg.MustParse(&args)

	appConfig, err := config.ParseConfigFile(args.ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	if !args.IsLive {
		log.Warn("Running in test mode, run with --live to submit real orders")
	}

	krakenAPI := krakenapi.New(appConfig.Key, appConfig.Secret)
	var notifiers []*notifier.Notifier

	if args.EmailCredentialsFile != "" {
		var gmailer notifier.Notifier = notifier.MustNewGMailer(args.EmailCredentialsFile, "me")
		notifiers = append(notifiers, &gmailer)
	} else {
		log.Warn("--email-credentials not set, email notifications are disabled")
	}

	if args.TelegramCredentialsFile != "" {
		var telegram notifier.Notifier = notifier.MustNewTelegramNotifier(args.TelegramCredentialsFile)
		notifiers = append(notifiers, &telegram)
	} else {
		log.Warn("--telegram-credentials not set, telegram notifications are disabled")
	}

	apiInstance := api.NewApi(*appConfig, args.IsLive, args.IsVerbose, krakenAPI)
	schedulerInstance := scheduler.NewScheduler(*appConfig, metrics.NewMetrics(), apiInstance, notifiers)

	retry.DefaultAttempts = 11
	retry.DefaultDelay = 60 * time.Second
	retry.DefaultDelayType = retry.BackOffDelay
	retry.DefaultOnRetry = func(n uint, err error) {
		subject := fmt.Sprintf("Retryable call failed (attempt %d)", n+1)
		log.Warnf("%s: %s", subject, err.Error())
		for _, n := range notifiers {
			_ = (*n).Send(subject, err.Error())
		}
	}

	go metrics.Start()
	schedulerInstance.Run()
}
