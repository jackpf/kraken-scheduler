package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Amount(t *testing.T) {
	order := NewOrder("testpair", 1000.0, 500.0)

	assert.Equal(t, order.Amount(), 0.5)
}
