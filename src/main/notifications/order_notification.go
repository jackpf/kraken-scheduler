package notifications

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"strings"
)

func NewOrderNotification(isLive bool, pair model.Pair, amount float64, orderPrice float64, assetPrice float64, transactionIds []string) Notification {
	return OrderNotification{isLive: isLive, pair: pair, amount: amount, orderPrice: orderPrice, assetPrice: assetPrice, transactionIds: transactionIds}
}

type OrderNotification struct {
	isLive         bool
	pair           model.Pair
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
	return fmt.Sprintf("kraken-scheduler%s: %s order submitted (%s)", testTag, n.pair.Name(), strings.Join(n.transactionIds[:], ", "))
}

func (n OrderNotification) Body() string {
	return fmt.Sprintf(`Transaction ID: %s.

Placed an order for %f%s, at a cost of %s%f.
Current asset price: %s = %s%f.

Purchase confirmation should arrive shortly, if not - check the application logs!`,
		strings.Join(n.transactionIds[:], ", "),
		n.amount,
		n.pair.First.Symbol,
		n.pair.Second.Symbol,
		n.orderPrice,
		n.pair.First.Name,
		n.pair.Second.Symbol,
		n.assetPrice)

}
