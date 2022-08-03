package notificationtemplates

import (
	"fmt"
	"strings"
)

func NewOrderNotification(isLive bool, pair string, amount float64, orderPrice float64, assetPrice float64, transactionIds []string) NotificationTemplate {
	return OrderNotification{isLive: isLive, pair: pair, amount: amount, orderPrice: orderPrice, assetPrice: assetPrice, transactionIds: transactionIds}
}

type OrderNotification struct {
	isLive         bool
	pair           string
	amount         float64
	orderPrice     float64
	assetPrice     float64
	transactionIds []string
}

func (o OrderNotification) Subject() string {
	testTag := ""
	if !o.isLive {
		testTag = " [TEST]"
	}
	return fmt.Sprintf("kraken-scheduler%s: %s order submitted (%s)", testTag, o.pair, strings.Join(o.transactionIds[:], ", "))
}

func (o OrderNotification) Body() string {
	return fmt.Sprintf(`Transaction ID: %s.

Placed an order for %f %s, at a cost of %f.
Current asset price: 1 * %s = %f.

Purchase confirmation should arrive shortly, if not - check the application logs!`,
		strings.Join(o.transactionIds[:], ", "),
		o.amount,
		o.pair,
		o.orderPrice,
		o.pair,
		o.assetPrice)

}
