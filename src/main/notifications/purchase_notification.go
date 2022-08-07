package notifications

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/util"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

func NewPurchaseNotification(pair model.Pair, amount float64, orderPrice float64, transactionId string, completedOrder krakenapi.Order) Notification {
	return PurchaseNotification{pair: pair, amount: amount, orderPrice: orderPrice, transactionId: transactionId, completedOrder: completedOrder}
}

type PurchaseNotification struct {
	pair           model.Pair
	amount         float64
	orderPrice     float64
	transactionId  string
	completedOrder krakenapi.Order // TODO Take individual values
}

func (n PurchaseNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: %s purchase successful (%s)", n.pair.Name(), n.transactionId)
}

func (n PurchaseNotification) Body() string {
	return fmt.Sprintf(`Transaction ID: %s.

Purchase of %s, at a cost of %s was successful.

Transaction summary: %+v`,
		n.transactionId,
		util.FormatAsset(n.pair.First, n.amount),
		util.FormatAsset(n.pair.Second, n.orderPrice),
		n.completedOrder)
}
