package model

import (
	"github.com/go-co-op/gocron"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
)

type TaskData struct {
	Jobs []struct {
		configmodel.Schedule
		*gocron.Job
	}
	Schedule       configmodel.Schedule
	Order          Order
	TransactionIds []string
	BalanceData    []apimodel.BalanceData
}

func (d TaskData) Job() *struct {
	configmodel.Schedule
	*gocron.Job
} {
	for _, job := range d.Jobs {
		if job.Schedule == d.Schedule {
			return &job
		}
	}

	return nil
}
