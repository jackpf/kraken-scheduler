package util

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAsset_Fiat(t *testing.T) {
	result := FormatAsset(configmodel.EUR, 123.456)

	assert.Equal(t, "€123.46", result)
}

func TestFormatAsset_Crypto(t *testing.T) {
	result := FormatAsset(configmodel.XBT, 123.4567891011)

	assert.Equal(t, "123.4567891₿", result)
}

func TestFormatAmount(t *testing.T) {
	result := FormatFloat(12.34567891011, 8)

	assert.Equal(t, "12.34567891", result)
}

func TestFormatCurrency(t *testing.T) {
	result := FormatCurrency(12.34567891011)

	assert.Equal(t, "12.35", result)
}
