package tasks

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/scheduler/model"
	"github.com/jackpf/kraken-scheduler/src/main/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrderTask_Run(t *testing.T) {
	api := new(testutil.MockApi)
	task := NewCreateOrderTask(api)
	taskData := model.TaskData{Schedule: configmodel.Schedule{Cron: "***", Amount: 123.0, Pair: configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}}}

	mockOrder := model.NewOrder(configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, 500.0, 123.0)

	api.On("CreateOrder", taskData.Schedule.Pair, taskData.Schedule.Amount).Return(&mockOrder, nil)

	err := task.Run(&taskData)

	assert.NoError(t, err)
	assert.Equal(t, mockOrder, taskData.Order)
	api.AssertExpectations(t)
}
