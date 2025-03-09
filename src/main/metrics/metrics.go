package metrics

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var assetLabels = []string{"asset", "asset_symbol"}
var currencyLabels = []string{"currency", "currency_symbol"}

func NewMetrics() Metrics {
	return &MetricsImpl{
		orderCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_orders_total",
			Help: "The total number of orders made",
		}, append(assetLabels, currencyLabels...)),
		purchaseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "kraken_scheduler_purchases_total",
			Help: "The total number of purchases made",
		}, append(assetLabels, currencyLabels...)),
		purchaseAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_purchase_amount",
			Help: "How much of an asset was purchased",
		}, append(assetLabels, currencyLabels...)),
		purchaseAmountHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_purchase_amount_history",
			Help: "How much of an asset was purchased",
		}, append(assetLabels, currencyLabels...)),
		spendAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_spend_amount",
			Help: "How much currency was spent on a purchase",
		}, append(assetLabels, currencyLabels...)),
		spendAmountHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "kraken_scheduler_spend_amount_history",
			Help: "How much currency was spent on a purchase",
		}, append(assetLabels, currencyLabels...)),
		assetBalanceAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_balance_amount",
			Help: "How much of an asset currently exists on the account",
		}, assetLabels),
		assetBalanceValueGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_asset_balance_value",
			Help: "How much value of an asset exists on the account",
		}, assetLabels),
		currencyBalanceAmountGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kraken_scheduler_currency_balance_amount",
			Help: "How much of a currency currently exists on the account",
		}, currencyLabels),
		errorsCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "kraken_scheduler_errors_total",
			Help: "The total number of errors during purchase process",
		}),
		retriesCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "kraken_scheduler_retries_total",
			Help: "The total number of retries during purchase process",
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
	orderCounter               *prometheus.CounterVec
	purchaseCounter            *prometheus.CounterVec
	purchaseAmountGauge        *prometheus.GaugeVec
	purchaseAmountHistogram    *prometheus.HistogramVec
	spendAmountGauge           *prometheus.GaugeVec
	spendAmountHistogram       *prometheus.HistogramVec
	assetBalanceAmountGauge    *prometheus.GaugeVec
	assetBalanceValueGauge     *prometheus.GaugeVec
	currencyBalanceAmountGauge *prometheus.GaugeVec
	errorsCounter              prometheus.Counter
	retriesCounter             prometheus.Counter
}

func pairToAssetLabels(pair model.Pair) []string {
	return []string{pair.First.Name, pair.First.Symbol}
}

func pairToCurrencyLabels(pair model.Pair) []string {
	return []string{pair.Second.Name, pair.Second.Symbol}
}

func pairToAssetAndCurrencyLabels(pair model.Pair) []string {
	return append(pairToAssetLabels(pair), pairToCurrencyLabels(pair)...)
}

func (m *MetricsImpl) LogOrder(pair model.Pair) {
	fmt.Printf("Logging order, first: %s, second: %s", pair.First.Name, pair.Second.Name)
	m.orderCounter.WithLabelValues(pairToAssetAndCurrencyLabels(pair)...).Inc()
}

func (m *MetricsImpl) LogPurchase(pair model.Pair, amount float64, fiatAmount float64, holdings float64, holdingsValue float64) {
	m.purchaseCounter.WithLabelValues(pairToAssetAndCurrencyLabels(pair)...).Inc()
	m.purchaseAmountGauge.WithLabelValues(pairToAssetAndCurrencyLabels(pair)...).Set(amount)
	m.purchaseAmountHistogram.WithLabelValues(pairToAssetAndCurrencyLabels(pair)...).Observe(amount)
	m.spendAmountGauge.WithLabelValues(pairToAssetAndCurrencyLabels(pair)...).Set(fiatAmount)
	m.spendAmountHistogram.WithLabelValues(pairToAssetAndCurrencyLabels(pair)...).Observe(fiatAmount)
	m.assetBalanceAmountGauge.WithLabelValues(pairToAssetLabels(pair)...).Set(holdings)
	m.assetBalanceValueGauge.WithLabelValues(pairToAssetLabels(pair)...).Set(holdingsValue)
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
