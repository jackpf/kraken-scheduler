package main

import (
	"github.com/jackpf/kraken-schedule/src/main/config"
	"github.com/jackpf/kraken-schedule/src/main/scheduler"
	log "github.com/sirupsen/logrus"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	"github.com/alexflint/go-arg"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

func main() {
	var args struct {
		Key             string `arg:"required" help:"Your Kraken API key"`
		Secret          string `arg:"required" help:"Your Kraken secret key"`
		ConfigFile      string `arg:"--config,required" help:"Schedule configuration file"`
		CredentialsFile string `arg:"--credentials" help:"Your google OAuth credentials.json file (optional)"`
		IsLive          bool   `arg:"--live" default:"false" help:"Set to true to execute real orders"`
	}
	arg.MustParse(&args)

	appConfig, err := config.ParseConfigFile(args.ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	if !args.IsLive {
		log.Warn("Running in test mode, run with `--live true` to submit real orders")
	}

	api := krakenapi.New(args.Key, args.Secret)
	var notifierInstance *notifier.Notifier
	if args.CredentialsFile != "" {
		var gmailer notifier.Notifier = notifier.MustNewGMailer(args.CredentialsFile, "me")
		notifierInstance = &gmailer
	} else {
		log.Warn("--credentials not set, notifications are disabled")
	}
	schedulerInstance := scheduler.NewScheduler(*appConfig, args.IsLive, api, notifierInstance)

	schedulerInstance.Run()
}
