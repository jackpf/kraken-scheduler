package model

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
)

type TaskData struct {
	Schedule       configmodel.Schedule
	Order          Order
	TransactionIds []string
}
