package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	"github.com/jackpf/kraken-scheduler/src/main/metrics"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	log "github.com/sirupsen/logrus"
)

const alertThreshold = 1.5

func NewCheckBalanceTask(api api.Api, metrics metrics.Metrics) CheckBalanceTask {
	return CheckBalanceTask{api: api, metrics: metrics}
}

type CheckBalanceTask struct {
	api     api.Api
	metrics metrics.Metrics
}

func (t CheckBalanceTask) Run(taskData *model.TaskData) error {
	var requests []apimodel.BalanceRequest

	for _, job := range taskData.Jobs {
		requests = append(requests, apimodel.BalanceRequest{Pair: job.Pair, Amount: job.Amount})
	}

	balanceInfo, err := t.api.CheckBalance(requests)
	if err != nil {
		return err
	}

	for _, balance := range balanceInfo {
		t.metrics.LogCurrencyBalance(balance.Asset, balance.Balance)
	}
	taskData.BalanceData = balanceInfo

	return nil
}

func (t CheckBalanceTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	var balanceNotifications []notifications.Notification

	for _, balance := range taskData.BalanceData {
		if balance.Balance/balance.NextPurchaseAmount < alertThreshold {
			log.Warnf("Low balance on account for %s", balance.Asset.Name)

			balanceNotifications = append(balanceNotifications, notifications.NewLowBalanceNotification(
				balance.Asset,
				balance.NextPurchaseAmount,
				balance.Balance,
			))
		}
	}

	return balanceNotifications, nil
}
