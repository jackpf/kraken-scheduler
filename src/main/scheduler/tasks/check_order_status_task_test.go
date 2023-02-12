package tasks

import (
	"fmt"
	krakenapi "github.com/beldur/kraken-go-api-client"
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/jackpf/kraken-scheduler/src/main/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckOrderStatusTask_Notifications(t *testing.T) {
	api := new(testutil.MockApi)
	task := NewCheckOrderStatusTask(api)
	taskData := model.TaskData{
		Schedule:       configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}},
		Order:          model.Order{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Price: 500.0, FiatAmount: 123.0},
		TransactionIds: []string{"1", "2"},
	}

	mockCompletedOrder1 := krakenapi.Order{TransactionID: "1"}
	mockCompletedOrder2 := krakenapi.Order{TransactionID: "2"}
	mockHoldings := 5000.0

	api.On("IsLive").Return(false)
	api.On("IsVerbose").Return(true)
	api.On("TransactionStatus", "1").Return(&mockCompletedOrder1, nil)
	api.On("TransactionStatus", "2").Return(&mockCompletedOrder2, nil)
	api.On("CheckHoldings", configmodel.XXBT).Return(&mockHoldings, nil)

	result, errs := task.Notifications(taskData)

	for _, err := range errs {
		assert.NoError(t, err)
	}

	assert.Equal(t, []notifications.Notification{
		notifications.NewPurchaseNotification(
			configmodel.Pair{configmodel.XXBT, configmodel.ZEUR},
			taskData.Order.Amount(),
			123.0,
			"1",
			mockCompletedOrder1,
			500.0,
			5000.0,
			true,
		),
		notifications.NewPurchaseNotification(
			configmodel.Pair{configmodel.XXBT, configmodel.ZEUR},
			taskData.Order.Amount(),
			123.0,
			"2",
			mockCompletedOrder2,
			500.0,
			5000.0,
			true,
		),
	}, result)
}

func TestCheckOrderStatusTask_Notifications_IfSomeFail(t *testing.T) {
	api := new(testutil.MockApi)
	task := NewCheckOrderStatusTask(api)
	taskData := model.TaskData{
		Schedule:       configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}},
		Order:          model.Order{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Price: 500.0, FiatAmount: 123.0},
		TransactionIds: []string{"1", "2", "3"},
	}

	mockCompletedOrder1 := krakenapi.Order{TransactionID: "1"}
	mockCompletedOrder2 := krakenapi.Order{}
	mockCompletedOrder3 := krakenapi.Order{TransactionID: "3"}
	mockHoldings := 5000.0

	api.On("IsLive").Return(false)
	api.On("IsVerbose").Return(true)
	api.On("TransactionStatus", "1").Return(&mockCompletedOrder1, nil)
	api.On("TransactionStatus", "2").Return(&mockCompletedOrder2, fmt.Errorf("mock error"))
	api.On("TransactionStatus", "3").Return(&mockCompletedOrder3, nil)
	api.On("CheckHoldings", configmodel.XXBT).Return(&mockHoldings, nil)

	result, errs := task.Notifications(taskData)

	assert.Equal(t, []notifications.Notification{
		notifications.NewPurchaseNotification(
			configmodel.Pair{configmodel.XXBT, configmodel.ZEUR},
			taskData.Order.Amount(),
			123.0,
			"1",
			mockCompletedOrder1,
			500,
			5000,
			true,
		),
		notifications.NewPurchaseNotification(
			configmodel.Pair{configmodel.XXBT, configmodel.ZEUR},
			taskData.Order.Amount(),
			123.0,
			"3",
			mockCompletedOrder3,
			500,
			5000,
			true,
		),
	}, result)

	assert.Len(t, errs, 1)
	assert.Errorf(t, errs[0], "mock error")
}
