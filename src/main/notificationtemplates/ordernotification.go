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

func (n OrderNotification) Subject() string {
	testTag := ""
	if !n.isLive {
		testTag = " [TEST]"
	}
	return fmt.Sprintf("kraken-scheduler%s: %s order submitted (%s)", testTag, n.pair, strings.Join(n.transactionIds[:], ", "))
}

func (n OrderNotification) Body() string {
	return fmt.Sprintf(`Transaction ID: %s.

Placed an order for %f %s, at a cost of %f.
Current asset price: 1 * %s = %f.

Purchase confirmation should arrive shortly, if not - check the application logs!`,
		strings.Join(n.transactionIds[:], ", "),
		n.amount,
		n.pair,
		n.orderPrice,
		n.pair,
		n.assetPrice)

}
