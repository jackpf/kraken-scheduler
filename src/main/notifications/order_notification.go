package notifications

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/util"
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

Placed an order for %s, at a cost of %s.
Current asset price: %s = %s.

Purchase confirmation should arrive shortly, if not - check the application logs!`,
		strings.Join(n.transactionIds[:], ", "),
		util.FormatAsset(n.pair.First, n.amount),
		util.FormatAsset(n.pair.Second, n.orderPrice),
		n.pair.First.Name,
		util.FormatAsset(n.pair.Second, n.assetPrice))
}
