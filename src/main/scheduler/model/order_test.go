package model

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Amount(t *testing.T) {
	order := NewOrder(configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, 1000.0, 500.0)

	assert.Equal(t, order.Amount(), 0.5)
}
