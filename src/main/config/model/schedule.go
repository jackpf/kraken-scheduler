package model

type Schedule struct {
	Cron   string  `json:"cron"`
	Pair   string  `json:"pair"`
	Amount float32 `json:"amount"`
}
