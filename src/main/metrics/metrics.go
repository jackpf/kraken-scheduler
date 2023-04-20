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
		}, []string{"currency", "asset"}),
		purchaseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_purchases_total",
			Help: "The total number of purchases made",
		}, []string{"currency", "asset"}),
		assetPurchaseAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_purchase_amount",
			Help: "How much of an asset was purchased",
		}, []string{"asset"}),
		currencySpendAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_currency_spend_amount",
			Help: "How much currency was spent on a purchase",
		}, []string{"currency"}),
		assetBalanceAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_balance_amount",
			Help: "How much of an asset currently exists on the account",
		}, []string{"asset"}),
		currencyBalanceAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_currency_balance_amount",
			Help: "How much of a currency currently exists on the account",
		}, []string{"currency"}),
		errorsCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_errors_total",
			Help: "The total number of errors during purchase process",
		}, []string{"currency", "asset"}),
	}
}

type Metrics interface {
	LogOrder(pair model.Pair)
	LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64)
	LogCurrencyBalance(asset model.Asset, holdings float64)
}

type MetricsImpl struct {
	orderCounter               *prometheus.CounterVec
	purchaseCounter            *prometheus.CounterVec
	assetPurchaseAmountGauge   *prometheus.GaugeVec
	currencySpendAmountGauge   *prometheus.GaugeVec
	assetBalanceAmountGauge    *prometheus.GaugeVec
	currencyBalanceAmountGauge *prometheus.GaugeVec
	errorsCounter              *prometheus.CounterVec
}

func (m *MetricsImpl) LogOrder(pair model.Pair) {
	m.orderCounter.WithLabelValues(pair.First.Name, pair.Second.Name).Inc()
}

func (m *MetricsImpl) LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64) {
	m.purchaseCounter.WithLabelValues(pair.First.Name, pair.Second.Name).Inc()
	m.assetPurchaseAmountGauge.WithLabelValues(pair.Second.Name).Set(amount)
	m.currencySpendAmountGauge.WithLabelValues(pair.First.Name).Set(fiatAmount)
	m.assetBalanceAmountGauge.WithLabelValues(pair.Second.Name).Set(holdings)
}

func (m *MetricsImpl) LogCurrencyBalance(asset model.Asset, holdings float64) {
	m.currencyBalanceAmountGauge.WithLabelValues(asset.Name).Set(holdings)
}
