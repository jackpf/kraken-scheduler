package model

type Schedule struct {
	Frequency string  `json:"frequency"`
	Day       int     `json:"day"`
	Time      string  `json:"time"`
	Pair      string  `json:"pair"`
	Amount    float32 `json:"amount"`
}
