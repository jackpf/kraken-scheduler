package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalanceRequest_Currency(t *testing.T) {
	balanceRequest := BalanceRequest{Pair: "XXBTZEUR", Amount: 123.0}

	assert.Equal(t, "ZEUR", balanceRequest.Currency())
}
