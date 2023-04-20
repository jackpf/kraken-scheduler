package tasks

import (
	"github.com/go-co-op/gocron"
	apimodel "github.com/jackpf/kraken-scheduler/src/main/api/model"
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/jackpf/kraken-scheduler/src/main/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckBalanceTask_Run(t *testing.T) {
	api := new(testutil.MockApi)
	metrics := new(testutil.MockMetrics)
	task := NewCheckBalanceTask(api, metrics)
	schedule := configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}}
	jobs := []struct {
		configmodel.Schedule
		*gocron.Job
	}{{schedule, nil}}
	taskData := model.TaskData{Jobs: jobs, Schedule: schedule}

	balanceInfo := []apimodel.BalanceData{{Asset: configmodel.ZEUR, NextPurchaseAmount: 123.0, Balance: 456.0}}

	api.On("CheckBalance", []apimodel.BalanceRequest{{Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, Amount: 123.0}}).Return(balanceInfo, nil)
	metrics.On("LogCurrencyBalance", configmodel.ZEUR, 456.0).Return()

	err := task.Run(&taskData)

	assert.NoError(t, err)
	assert.Equal(t, balanceInfo, taskData.BalanceData)
}

func TestCheckBalanceTask_Notifications(t *testing.T) {
	api := new(testutil.MockApi)
	metrics := new(testutil.MockMetrics)
	task := NewCheckBalanceTask(api, metrics)
	taskData := model.TaskData{BalanceData: []apimodel.BalanceData{
		{Asset: configmodel.ZUSD, NextPurchaseAmount: 100.0, Balance: 150.0},
		{Asset: configmodel.ZEUR, NextPurchaseAmount: 100.0, Balance: 149.99},
	}}

	balanceNotifications, errs := task.Notifications(taskData)

	for _, err := range errs {
		assert.NoError(t, err)
	}

	assert.Len(t, balanceNotifications, 1)
	assert.Equal(t, notifications.NewLowBalanceNotification(configmodel.ZEUR, 100.0, 149.99), balanceNotifications[0])
}
