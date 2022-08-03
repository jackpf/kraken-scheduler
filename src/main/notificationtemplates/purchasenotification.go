package notificationtemplates

import (
	"fmt"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

func NewPurchaseNotification(pair string, amount float64, orderPrice float64, transactionId string, completedOrder krakenapi.Order) NotificationTemplate {
	return PurchaseNotification{pair: pair, amount: amount, orderPrice: orderPrice, transactionId: transactionId, completedOrder: completedOrder}
}

type PurchaseNotification struct {
	pair           string
	amount         float64
	orderPrice     float64
	transactionId  string
	completedOrder krakenapi.Order // TODO Take individual values
}

func (p PurchaseNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: %s purchase successful (%s)", p.pair, p.transactionId)
}

func (p PurchaseNotification) Body() string {
	return fmt.Sprintf(`Transaction ID: %s.

Purchase of %f %s, at a cost of %f was successful.

Transaction summary: %+v`,
		p.transactionId,
		p.amount,
		p.pair,
		p.orderPrice,
		p.completedOrder)
}
