package main

import (
	"fmt"
	"log"

	"github.com/jackpf/kraken-schedule/src/main/notifier"

	"github.com/alexflint/go-arg"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/jackpf/kraken-schedule/src/main/calculator"
	"github.com/jackpf/kraken-schedule/src/main/config"
)

func main() {
	var args struct {
		Key              string `arg:"required"`
		Secret           string `arg:"required"`
		ScheduleFileName string `arg:"-f,--file,required"`
	}
	arg.MustParse(&args)

	api := krakenapi.New(args.Key, args.Secret)
	result, err := api.Balance()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current balance: %+v\n", result)

	appConfig, err := config.ParseConfigFile(args.ScheduleFileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Schedule: %+v\n", appConfig)

	currencyCalculator := calculator.NewCurrencyCalculator(api)
	amount, err := currencyCalculator.AmountFor(appConfig.Schedules[0].Pair, appConfig.Schedules[0].Amount)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Buying %+v btc for €%+v\n", *amount, appConfig.Schedules[0].Amount)

	mailer, err := notifier.NewGMailer("credentials.json", "me")

	if err != nil {
		log.Fatal(err)
	}

	err = mailer.Send(appConfig.NotifyEmailAddress, "kraken-scheduler: Purchase", fmt.Sprintf("Buying %+v btc for €%+v\n", *amount, appConfig.Schedules[0].Amount))

	if err != nil {
		log.Fatal(err)
	}
}
