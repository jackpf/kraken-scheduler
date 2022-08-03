package main

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	"github.com/jackpf/kraken-scheduler/src/main/config"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler"
	log "github.com/sirupsen/logrus"

	"github.com/jackpf/kraken-scheduler/src/main/notifier"

	"github.com/alexflint/go-arg"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

func main() {
	var args struct {
		Key                     string `arg:"required" help:"Your Kraken API key"`
		Secret                  string `arg:"required" help:"Your Kraken secret key"`
		ConfigFile              string `arg:"--config,required" help:"Schedule configuration file"`
		EmailCredentialsFile    string `arg:"--email-credentials" help:"Your google OAuth email-credentials.json file (optional)"`
		TelegramCredentialsFile string `arg:"--telegram-credentials" help:"Your telegram ChatID and Token telegram-credentials.json file (optional)"`
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

	krakenAPI := krakenapi.New(args.Key, args.Secret)
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

	apiInstance := api.NewApi(*appConfig, args.IsLive, krakenAPI)
	schedulerInstance := scheduler.NewScheduler(*appConfig, apiInstance, notifiers)

	schedulerInstance.Run()
}
