package model

type Schedule struct {
	Cron   string  `json:"cron"`
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}
