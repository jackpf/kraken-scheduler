package metrics

import (
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

var metrics = NewMetrics().(*MetricsImpl)
var pair = configmodel.Pair{configmodel.XXBT, configmodel.ZEUR}
var expectedAssetLabels = []string{"XXBT", "₿"}
var expectedCurrencyLabels = []string{"ZEUR", "€"}
var expectedAssetAndCurrencyLabels = append(expectedAssetLabels, expectedCurrencyLabels...)

func TestMetrics_LogOrder(t *testing.T) {
	metrics.LogOrder(pair)

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.orderCounter))
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.orderCounter.WithLabelValues(expectedAssetAndCurrencyLabels...)))
}

func TestMetrics_LogPurchase(t *testing.T) {
	amount := float64(300)
	fiatAmount := float64(50)
	holdings := float64(600)
	holdingsValue := float64(100)
	metrics.LogPurchase(pair, amount, fiatAmount, holdings, holdingsValue)

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.purchaseCounter))
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.purchaseCounter.WithLabelValues(expectedAssetAndCurrencyLabels...)))

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.purchaseAmountGauge))
	assert.Equal(t, amount, testutil.ToFloat64(metrics.purchaseAmountGauge.WithLabelValues(expectedAssetAndCurrencyLabels...)))

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.purchaseAmountHistogram))
	// TODO Not sure how to test histogram values...

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.spendAmountGauge))
	assert.Equal(t, amount, testutil.ToFloat64(metrics.purchaseAmountGauge.WithLabelValues(expectedAssetAndCurrencyLabels...)))

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.spendAmountHistogram))
	// TODO Not sure how to test histogram values...

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.assetBalanceAmountGauge))
	assert.Equal(t, holdings, testutil.ToFloat64(metrics.assetBalanceAmountGauge.WithLabelValues(expectedAssetLabels...)))

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.assetBalanceValueGauge))
	assert.Equal(t, holdingsValue, testutil.ToFloat64(metrics.assetBalanceValueGauge.WithLabelValues(expectedAssetLabels...)))
}

func TestMetrics_LogCurrencyBalance(t *testing.T) {
	holdings := float64(600)
	metrics.LogCurrencyBalance(pair.First, holdings)

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.currencyBalanceAmountGauge))
	assert.Equal(t, holdings, testutil.ToFloat64(metrics.currencyBalanceAmountGauge.WithLabelValues(expectedAssetLabels...)))
}

func TestMetrics_LogError(t *testing.T) {
	metrics.LogError()

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.errorsCounter))
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.errorsCounter))
}

func TestMetrics_LogRetry(t *testing.T) {
	metrics.LogRetry()

	assert.Equal(t, 1, testutil.CollectAndCount(metrics.retriesCounter))
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.retriesCounter))
}
