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
	task := NewCheckBalanceTask(api)
	schedule := configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: "XXBTZEUR"}
	jobs := []struct {
		configmodel.Schedule
		*gocron.Job
	}{{schedule, nil}}
	taskData := model.TaskData{Jobs: jobs, Schedule: schedule}

	balanceInfo := []apimodel.BalanceData{{Currency: "ZEUR", NextPurchaseAmount: 123.0, Balance: 456.0}}

	api.On("CheckBalance", []apimodel.BalanceRequest{{Pair: "XXBTZEUR", Amount: 123.0}}).Return(balanceInfo, nil)

	err := task.Run(&taskData)

	assert.NoError(t, err)
	assert.Equal(t, balanceInfo, taskData.BalanceData)
}

func TestCheckBalanceTask_Notifications(t *testing.T) {
	api := new(testutil.MockApi)
	task := NewCheckBalanceTask(api)
	taskData := model.TaskData{BalanceData: []apimodel.BalanceData{
		{Currency: "ZUSD", NextPurchaseAmount: 100.0, Balance: 150.0},
		{Currency: "ZEUR", NextPurchaseAmount: 100.0, Balance: 149.99},
	}}

	balanceNotifications, errs := task.Notifications(taskData)

	for _, err := range errs {
		assert.NoError(t, err)
	}

	assert.Len(t, balanceNotifications, 1)
	assert.Equal(t, notifications.NewLowBalanceNotification("ZEUR", 100.0, 149.99), balanceNotifications[0])
}
