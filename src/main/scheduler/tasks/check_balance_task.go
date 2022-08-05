package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	log "github.com/sirupsen/logrus"
)

const alertThreshold = 1.5

func NewCheckBalanceTask(api api.Api) CheckBalanceTask {
	return CheckBalanceTask{api: api}
}

type CheckBalanceTask struct {
	api api.Api
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

	taskData.BalanceData = balanceInfo

	return nil
}

func (t CheckBalanceTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	var balanceNotifications []notifications.Notification

	for _, balance := range taskData.BalanceData {
		if balance.Balance/balance.NextPurchaseAmount < alertThreshold {
			log.Warnf("Low balance on account for %s", balance.Currency)

			balanceNotifications = append(balanceNotifications, notifications.NewLowBalanceNotification(
				balance.Currency,
				balance.NextPurchaseAmount,
				balance.Balance,
			))
		}
	}

	return balanceNotifications, nil
}
