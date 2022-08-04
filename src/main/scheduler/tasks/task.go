package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
)

type Task interface {
	Run(taskData *model.TaskData) (*model.TaskData, error)
	Notifications(taskData *model.TaskData) ([]notifications.Notification, []error)
}
