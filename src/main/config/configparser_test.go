package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigValid(t *testing.T) {
	configuration := []byte(`{
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

	assert.Len(t, config.Schedules, 1)
	assert.Equal(t, "* * * * *", config.Schedules[0].Cron)
	assert.Equal(t, "XXBTZEUR", config.Schedules[0].Pair)
	assert.Equal(t, 40.00, config.Schedules[0].Amount)
}
