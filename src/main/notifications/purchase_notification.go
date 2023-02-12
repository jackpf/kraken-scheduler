package notifications

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/util"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

func NewPurchaseNotification(pair model.Pair, amount float64, orderPrice float64, transactionId string, completedOrder krakenapi.Order, assetPrice float64, holdings float64, verbose bool) Notification {
	return PurchaseNotification{pair: pair, amount: amount, orderPrice: orderPrice, transactionId: transactionId, completedOrder: completedOrder, assetPrice: assetPrice, holdings: holdings, verbose: verbose}
}

type PurchaseNotification struct {
	pair           model.Pair
	amount         float64
	orderPrice     float64
	transactionId  string
	completedOrder krakenapi.Order // TODO Take individual values
	assetPrice     float64
	holdings       float64
	verbose        bool
}

func (n PurchaseNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: %s purchase successful (%s)", n.pair.Name(), n.transactionId)
}

func (n PurchaseNotification) Body() string {
	transactionSummary := ""
	if n.verbose {
		transactionSummary = fmt.Sprintf("\n\nTransaction summary: %+v", n.completedOrder)
	}

	return fmt.Sprintf(`Transaction ID: %s.

Purchase of %s for %s was successful.

Current holdings: %s
Holdings value: %s.%s`,
		n.transactionId,
		util.FormatAsset(n.pair.First, n.amount),
		util.FormatAsset(n.pair.Second, n.orderPrice),
		util.FormatAsset(n.pair.First, n.holdings),
		util.FormatAsset(n.pair.Second, n.holdings*n.assetPrice),
		transactionSummary)
}
