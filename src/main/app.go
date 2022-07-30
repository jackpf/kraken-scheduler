package main

import (
	"log"

	"github.com/jackpf/kraken-schedule/src/main/config"
	scheduler "github.com/jackpf/kraken-schedule/src/main/scheduler"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	"github.com/alexflint/go-arg"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

func main() {
	var args struct {
		Key        string `arg:"required"`
		Secret     string `arg:"required"`
		ConfigFile string `arg:"-f,--file,required"`
	}
	arg.MustParse(&args)

	appConfig, err := config.ParseConfigFile(args.ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	api := krakenapi.New(args.Key, args.Secret)
	notifier := notifier.MustNewGMailer("credentials.json", "me")
	scheduler := scheduler.NewScheduler(*appConfig, api, notifier)

	scheduler.Run()
}
