package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
)

func NewCreateOrderTask(api api.Api) CreateOrderTask {
	return CreateOrderTask{api: api}
}

type CreateOrderTask struct {
	api api.Api
}

func (t CreateOrderTask) Run(taskData *model.TaskData) error {
	order, err := t.api.CreateOrder(taskData.Schedule.Pair, taskData.Schedule.Amount)
	if err != nil {
		return err
	}

	taskData.Order = *order
	return nil
}

func (t CreateOrderTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	return nil, nil
}
