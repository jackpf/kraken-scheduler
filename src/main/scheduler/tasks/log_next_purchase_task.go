package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	log "github.com/sirupsen/logrus"
)

func NewLogNextPurchaseTask() LogNextPurchaseTask {
	return LogNextPurchaseTask{}
}

type LogNextPurchaseTask struct {
}

func (t LogNextPurchaseTask) Run(taskData *model.TaskData) error {
	if taskData.Job() != nil {
		log.Infof("Next purchase for %s will occur at %+v", taskData.Job().Pair, taskData.Job().NextRun())
	}

	return nil
}

func (t LogNextPurchaseTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	return nil, nil
}
