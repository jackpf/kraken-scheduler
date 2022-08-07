package model

type Schedule struct {
	Cron   string  `json:"cron"`
	Pair   Pair    `json:"pair"`
	Amount float64 `json:"amount"`
}
