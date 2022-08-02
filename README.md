# Kraken Scheduler

<img src="./doc/img/kraken-logo.png" alt="drawing" width="200" />

Build status:

![build status](https://github.com/jackpf/kraken-scheduler/actions/workflows/go.yml/badge.svg)


Tired of manually buying crypto every month? Every week? Every day?!

This application creates automated buy orders for cryptocurrencies on [Kraken](https://www.kraken.com/)
based on your configuration, with email alerts on order placements and status (currently requires GMail).

Disclaimer: this application isn't affiliated with Kraken in any way, and I take no responsibility
for incorrectly placed orders.

## Prerequisites

- You must create a Kraken API key to run the scheduler, see https://support.kraken.com/hc/en-us/articles/360000919966-How-to-generate-an-API-key-pair-
- Required permissions are `Query funds`, `Query open orders & trades`, `Query closed orders & trades` and `Create and modify orders` for the application to run correctly
- In order to receive email notifications, you must create your own GMail OAuth credentials, see https://developers.google.com/identity/protocols/oauth2

## Installation from binaries

Check the [releases](https://github.com/jackpf/kraken-scheduler/releases), and download the binary relating
to your operating system.

## Building

1. Ensure you have [go](https://go.dev/) installed (at least version 1.18)
2. Run `make build`
3. Executable is created in `./target/kraken-schedule`

## Installation from source (Linux / OSX only)

1. Run `make install` (binary is copied to `/usr/local/bin`)

## Schedule configuration

The application needs a JSON configuration file to run.

Create `config.json` in the directory of your choice (eg. `$HOME/.kraken-schedule/config.json`).

Example configuration:

```json
{
  "notify": "your-email@address.com", // optional - set this to receive email notifications
  "schedules": [
    {
      "cron":  "00 12 * * 1",
      "pair": "XXBTZEUR",
      "amount": 100.00
    },
    {
      "cron":  "00 12 1 * *",
      "pair": "ADAEUR",
      "amount": 50.00
    }
  ]
}
```

This example will order €100 of bitcoin every week, and €50 of ADA every month.

Here is a detailed explanation of each schedule parameter:

| Parameter 	| Description                                                                                                                                                      	| Valid Values                                                                                                                                                                                     	|
|-----------	|------------------------------------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|
| cron      	| A crontab configuration describing when to create orders for this schedule (use [crontab.guru](https://crontab.guru/) for help)                                  	| A valid crontab string                                                                                                                                                                           	|
| pair      	| Crypto & fiat pair to purchase, eg. XXBTZEUR to purchase bitcoin in euros, ADAEUR to purchase ADA in USD.ADAUSD                                                  	| Check the  [kraken go client](https://github.com/beldur/kraken-go-api-client/blob/master/types.go#L14-L88)  for all supported pairs                                                              	|
| amount    	| Amount of crypto to purchase, in fiat. This depends on the `pair` - if `XXBTZEUR` for example, `amount` will be in euros, if `XXBTZUSD`, `amount` will be in USD 	| Any float (check [kraken minimum order amounts](https://support.kraken.com/hc/en-us/articles/205893708-Minimum-order-size-volume-for-trading) - if the amount is too small your order will fail) 	|

## Running

Run with:

```shell
kraken-schedule --key KEY --secret SECRET --config CONFIG [--credentials CREDENTIALS] [--live]
```

Note that by default the application runs in test mode, and doesn't create real orders.

This is useful to validate that you've configured things correctly, and the purchase amounts are correct.

To place real orders, you must pass `--live` when running.

Run `kraken-schedule --help` for a description of all arguments.
