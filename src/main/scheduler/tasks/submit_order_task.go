package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	"github.com/jackpf/kraken-scheduler/src/main/metrics"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/jackpf/kraken-scheduler/src/main/util"
	log "github.com/sirupsen/logrus"
	"strings"
)

func NewSubmitOrderTask(api api.Api, metrics metrics.Metrics) SubmitOrderTask {
	return SubmitOrderTask{api: api, metrics: metrics}
}

type SubmitOrderTask struct {
	api     api.Api
	metrics metrics.Metrics
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
	t.metrics.LogOrder(taskData.Order.Pair)

	taskData.TransactionIds = transactionIds
	return nil
}

func (t SubmitOrderTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	if !t.api.IsVerbose() {
		return nil, nil
	}

	return []notifications.Notification{notifications.NewOrderNotification(
		t.api.IsLive(),
		taskData.Order.Pair,
		taskData.Order.Amount(),
		taskData.Order.FiatAmount,
		taskData.Order.Price,
		taskData.TransactionIds,
	)}, nil
}
