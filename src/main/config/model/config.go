package model

type Config struct {
	Key       string     `json:"key"`
	Secret    string     `json:"secret"`
	Schedules []Schedule `json:"schedules"`
}
