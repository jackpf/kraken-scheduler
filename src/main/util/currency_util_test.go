package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAmount(t *testing.T) {
	result := FormatFloat(12.34567891011, 8)

	assert.Equal(t, "12.34567891", result)
}

func TestFormatCurrency(t *testing.T) {
	result := FormatCurrency(12.34567891011)

	assert.Equal(t, "12.35", result)
}
