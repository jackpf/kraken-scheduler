package tasks

import (
	"github.com/jackpf/kraken-scheduler/src/main/api"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewCheckOrderStatusTask(api api.Api) CheckOrderStatusTask {
	return CheckOrderStatusTask{api: api}
}

type CheckOrderStatusTask struct {
	api api.Api
}

func (t CheckOrderStatusTask) Run(taskData *model.TaskData) error {
	return nil
}

func (t CheckOrderStatusTask) Notifications(taskData model.TaskData) ([]notifications.Notification, []error) {
	var notificationsList []notifications.Notification
	var errs []error

	for _, transactionId := range taskData.TransactionIds {
		for { // TODO perform in background & have max attempts
			completedOrder, err := t.api.TransactionStatus(transactionId)

			if err != nil {
				errs = append(errs, err)
				break
			}

			if completedOrder != nil {
				log.Infof("Order %s was successfully completed", transactionId)

				balanceInfo, err := t.api.CheckBalance([]apimodel.BalanceRequest{{Pair: taskData.Order.Pair, Amount: taskData.Order.Amount()}})
				if err != nil {
					errs = append(errs, err)
					break
				}

				notification := notifications.NewPurchaseNotification(
					taskData.Order.Pair,
					taskData.Order.Amount(),
					taskData.Order.FiatAmount,
					transactionId,
					*completedOrder,
					balanceInfo[0],
					t.api.IsVerbose(),
				)

				notificationsList = append(notificationsList, notification)

				break
			} else {
				log.Infof("Order %s is pending...", transactionId)
				time.Sleep(1 * time.Second)
			}
		}
	}

	return notificationsList, errs
}
