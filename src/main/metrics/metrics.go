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
		}, []string{"asset", "asset_symbol", "currency", "currency_symbol"}),
		purchaseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_purchases_total",
			Help: "The total number of purchases made",
		}, []string{"asset", "asset_symbol", "currency", "currency_symbol"}),
		assetPurchaseAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_purchase_amount",
			Help: "How much of an asset was purchased",
		}, []string{"asset", "asset_symbol"}),
		assetPurchaseAmountHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_asset_purchase_amount_history",
			Help: "How much of an asset was purchased",
		}, []string{"asset", "asset_symbol"}),
		currencySpendAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_currency_spend_amount",
			Help: "How much currency was spent on a purchase",
		}, []string{"currency", "currency_symbol"}),
		currencySpendAmountHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_currency_spend_amount_history",
			Help: "How much currency was spent on a purchase",
		}, []string{"currency", "currency_symbol"}),
		assetBalanceAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_balance_amount",
			Help: "How much of an asset currently exists on the account",
		}, []string{"asset", "asset_symbol"}),
		assetBalanceValueGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_balance_value",
			Help: "How much value of an asset exists on the account",
		}, []string{"asset", "asset_symbol"}),
		currencyBalanceAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_currency_balance_amount",
			Help: "How much of a currency currently exists on the account",
		}, []string{"currency", "currency_symbol"}),
		errorsCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "kraken_scheduler_errors_total",
			Help: "The total number of errors during purchase process",
		}),
		retriesCounter: promauto.NewCounter(prometheus.CounterOpts{
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
	LogRetry()
}

type MetricsImpl struct {
	orderCounter                 *prometheus.CounterVec
	purchaseCounter              *prometheus.CounterVec
	assetPurchaseAmountGauge     *prometheus.GaugeVec
	assetPurchaseAmountHistogram *prometheus.HistogramVec
	currencySpendAmountGauge     *prometheus.GaugeVec
	currencySpendAmountHistogram *prometheus.HistogramVec
	assetBalanceAmountGauge      *prometheus.GaugeVec
	assetBalanceValueGauge       *prometheus.GaugeVec
	currencyBalanceAmountGauge   *prometheus.GaugeVec
	errorsCounter                prometheus.Counter
	retriesCounter               prometheus.Counter
}

func (m *MetricsImpl) LogOrder(pair model.Pair) {
	m.orderCounter.WithLabelValues(pair.First.Name, pair.First.Symbol, pair.Second.Name, pair.Second.Symbol).Inc()
}

func (m *MetricsImpl) LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64, holdingsValue float64) {
	m.purchaseCounter.WithLabelValues(pair.First.Name, pair.First.Symbol, pair.Second.Name, pair.Second.Symbol).Inc()
	m.assetPurchaseAmountGauge.WithLabelValues(pair.First.Name, pair.First.Symbol).Set(amount)
	m.assetPurchaseAmountHistogram.WithLabelValues(pair.First.Name, pair.First.Symbol).Observe(amount)
	m.currencySpendAmountGauge.WithLabelValues(pair.Second.Name, pair.Second.Symbol).Set(fiatAmount)
	m.currencySpendAmountHistogram.WithLabelValues(pair.Second.Name, pair.Second.Symbol).Observe(fiatAmount)
	m.assetBalanceAmountGauge.WithLabelValues(pair.First.Name, pair.First.Symbol).Set(holdings)
	m.assetBalanceValueGauge.WithLabelValues(pair.First.Name, pair.First.Symbol).Set(holdingsValue)
}

func (m *MetricsImpl) LogCurrencyBalance(asset model.Asset, holdings float64) {
	m.currencyBalanceAmountGauge.WithLabelValues(asset.Name, asset.Symbol).Set(holdings)
}

func (m *MetricsImpl) LogError() {
	m.errorsCounter.Inc()
}

func (m *MetricsImpl) LogRetry() {
	m.retriesCounter.Inc()
}
