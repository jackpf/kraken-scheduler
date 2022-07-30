package model

type Config struct {
	NotifyEmailAddress string     `json:"notify"`
	Schedules          []Schedule `json:"schedules"`
}
