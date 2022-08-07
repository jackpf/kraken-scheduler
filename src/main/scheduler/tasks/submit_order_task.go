package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/jackpf/kraken-scheduler/src/main/util"
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

func (t SubmitOrderTask) Run(taskData *model.TaskData) error {
	log.Infof("[%s] Ordering %s for %s (1%s = %s)...",
		t.liveLogTag(),
		util.FormatAsset(taskData.Order.Pair.First, taskData.Order.Amount()),
		util.FormatAsset(taskData.Order.Pair.Second, taskData.Order.FiatAmount),
		taskData.Order.Pair.First.Symbol,
		util.FormatAsset(taskData.Order.Pair.Second, taskData.Order.Price))

	transactionIds, err := t.api.SubmitOrder(taskData.Order)
	if err != nil {
		return err
	}

	transactionIdsString := strings.Join(transactionIds[:], ", ")
	if !t.api.IsLive() {
		transactionIdsString = "<no transaction IDs for test orders>"
	}

	log.Infof("[%s] Order placed: %s", t.liveLogTag(), transactionIdsString)

	taskData.TransactionIds = transactionIds
	return nil
}

func (t SubmitOrderTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	return []notifications.Notification{notifications.NewOrderNotification(
		t.api.IsLive(),
		taskData.Order.Pair,
		taskData.Order.Amount(),
		taskData.Order.FiatAmount,
		taskData.Order.Price,
		taskData.TransactionIds,
	)}, nil
}
