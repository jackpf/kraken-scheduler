package notifications

import (
	"fmt"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

func NewPurchaseNotification(pair string, amount float64, orderPrice float64, transactionId string, completedOrder krakenapi.Order) Notification {
	return PurchaseNotification{pair: pair, amount: amount, orderPrice: orderPrice, transactionId: transactionId, completedOrder: completedOrder}
}

type PurchaseNotification struct {
	pair           string
	amount         float64
	orderPrice     float64
	transactionId  string
	completedOrder krakenapi.Order // TODO Take individual values
}

func (n PurchaseNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: %s purchase successful (%s)", n.pair, n.transactionId)
}

func (n PurchaseNotification) Body() string {
	return fmt.Sprintf(`Transaction ID: %s.

Purchase of %f %s, at a cost of %f was successful.

Transaction summary: %+v`,
		n.transactionId,
		n.amount,
		n.pair,
		n.orderPrice,
		n.completedOrder)
}
