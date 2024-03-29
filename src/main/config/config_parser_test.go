package config

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigValid(t *testing.T) {
	configuration := []byte(`{
  "key": "api-key",
  "secret": "api-secret",
  "schedules": [
    {
      "cron":  "* * * * *",
      "pair": "XXBTZEUR",
      "amount": 40.00
    }
  ]
}`)

	config, err := ParseConfig(configuration)

	assert.NoError(t, err)

	assert.Equal(t, "api-key", config.Key)
	assert.Equal(t, "api-secret", config.Secret)
	assert.Len(t, config.Schedules, 1)
	assert.Equal(t, "* * * * *", config.Schedules[0].Cron)
	assert.Equal(t, configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}, config.Schedules[0].Pair)
	assert.Equal(t, 40.00, config.Schedules[0].Amount)
}
