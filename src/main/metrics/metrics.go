package metrics

import (
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewMetrics() Metrics {
	return &MetricsImpl{
		orderCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_orders_total",
			Help: "The total number of orders made",
		}, []string{"asset", "currency"}),
		purchaseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_purchases_total",
			Help: "The total number of purchases made",
		}, []string{"asset", "currency"}),
		assetPurchaseAmountGauge: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_asset_purchase_amount",
			Help: "How much of an asset was purchased",
		}, []string{"asset"}),
		currencySpendAmountGauge: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_currency_spend_amount",
			Help: "How much currency was spent on a purchase",
		}, []string{"currency"}),
		assetBalanceAmountGauge: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_asset_balance_amount",
			Help: "How much of an asset currently exists on the account",
		}, []string{"asset"}),
		assetBalanceValueGauge: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_asset_balance_value",
			Help: "How much value of an asset exists on the account",
		}, []string{"asset"}),
		currencyBalanceAmountGauge: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_currency_balance_amount",
			Help: "How much of a currency currently exists on the account",
		}, []string{"currency"}),
		errorsCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "kraken_scheduler_errors_total",
			Help: "The total number of errors during purchase process",
		}),
	}
}

type Metrics interface {
	LogOrder(pair model.Pair)
	LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64, holdingsValue float64)
	LogCurrencyBalance(asset model.Asset, holdings float64)
	LogError()
}

type MetricsImpl struct {
	orderCounter               *prometheus.CounterVec
	purchaseCounter            *prometheus.CounterVec
	assetPurchaseAmountGauge   *prometheus.HistogramVec
	currencySpendAmountGauge   *prometheus.HistogramVec
	assetBalanceAmountGauge    *prometheus.HistogramVec
	assetBalanceValueGauge     *prometheus.HistogramVec
	currencyBalanceAmountGauge *prometheus.HistogramVec
	errorsCounter              prometheus.Counter
}

func (m *MetricsImpl) LogOrder(pair model.Pair) {
	m.orderCounter.WithLabelValues(pair.First.Name, pair.Second.Name).Inc()
}

func (m *MetricsImpl) LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64, holdingsValue float64) {
	m.purchaseCounter.WithLabelValues(pair.First.Name, pair.Second.Name).Inc()
	m.assetPurchaseAmountGauge.WithLabelValues(pair.First.Name).Observe(amount)
	m.currencySpendAmountGauge.WithLabelValues(pair.Second.Name).Observe(fiatAmount)
	m.assetBalanceAmountGauge.WithLabelValues(pair.First.Name).Observe(holdings)
	m.assetBalanceValueGauge.WithLabelValues(pair.First.Name).Observe(holdings)
}

func (m *MetricsImpl) LogCurrencyBalance(asset model.Asset, holdings float64) {
	m.currencyBalanceAmountGauge.WithLabelValues(asset.Name).Observe(holdings)
}

func (m *MetricsImpl) LogError() {
	m.errorsCounter.Inc()
}
