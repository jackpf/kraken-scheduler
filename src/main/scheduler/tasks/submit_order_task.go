package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	"github.com/jackpf/kraken-scheduler/src/main/notificationtemplates"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

func NewSubmitOrderTask(api api.Api) SubmitOrderTask {
	return SubmitOrderTask{api: api}
}

type SubmitOrderTask struct {
	api api.Api
}

func (t SubmitOrderTask) liveLogTag() string {
	if t.api.IsLive() {
		return "LIVE"
	}
	return "TEST"
}

func (t SubmitOrderTask) Run(taskData *model.TaskData) (*model.TaskData, error) {
	log.Infof("[%s] Ordering %s %s for %+v (%s = %f)...",
		t.liveLogTag(),
		t.api.FormatAmount(taskData.Order.Amount()),
		taskData.Order.Pair,
		taskData.Order.FiatAmount,
		taskData.Order.Pair,
		taskData.Order.Price)

	transactionIds, err := t.api.SubmitOrder(taskData.Order)
	if err != nil {
		return nil, err
	}

	transactionIdsString := strings.Join(transactionIds[:], ", ")
	if !t.api.IsLive() {
		transactionIdsString = "<no transaction IDs for test orders>"
	}

	log.Infof("[%s] Order placed: %s", t.liveLogTag(), transactionIdsString)

	taskData.TransactionIds = transactionIds

	return taskData, nil
}

func (t SubmitOrderTask) Notifications(taskData *model.TaskData) ([]notificationtemplates.NotificationTemplate, error) {
	return []notificationtemplates.NotificationTemplate{notificationtemplates.NewOrderNotification(
		t.api.IsLive(),
		taskData.Order.Pair,
		taskData.Order.Amount(),
		taskData.Order.FiatAmount,
		taskData.Order.Price,
		taskData.TransactionIds,
	)}, nil
}
