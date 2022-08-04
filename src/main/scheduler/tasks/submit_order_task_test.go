package tasks

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/notifications"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/jackpf/kraken-scheduler/src/main/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitOrderTask_Run(t *testing.T) {
	api := new(testutil.MockApi)
	task := NewSubmitOrderTask(api)
	taskData := model.TaskData{
		Schedule: configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: "mock-pair"},
		Order:    model.Order{Pair: "mock-pair", Price: 500.0, FiatAmount: 123.0},
	}

	mockTransactionIds := []string{"1", "2"}

	api.On("IsLive").Return(false)
	api.On("SubmitOrder", taskData.Order).Return(mockTransactionIds, nil)

	err := task.Run(&taskData)

	assert.NoError(t, err)
	assert.Equal(t, mockTransactionIds, taskData.TransactionIds)
}

func TestSubmitOrderTask_Notifications(t *testing.T) {
	api := new(testutil.MockApi)
	task := NewSubmitOrderTask(api)
	taskData := model.TaskData{
		Schedule:       configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: "mock-pair"},
		Order:          model.Order{Pair: "mock-pair", Price: 500.0, FiatAmount: 123.0},
		TransactionIds: []string{"1", "2"},
	}

	api.On("IsLive").Return(false)

	result, errs := task.Notifications(taskData)

	for _, err := range errs {
		assert.NoError(t, err)
	}
	assert.Equal(t, []notifications.Notification{notifications.NewOrderNotification(
		false,
		"mock-pair",
		taskData.Order.Amount(),
		123.0,
		500.0,
		[]string{"1", "2"},
	)}, result)
}
