package notifications

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	"github.com/jackpf/kraken-scheduler/src/main/util"
)

func NewLowBalanceNotification(asset model.Asset, fiatAmount float64, balanceAmount float64) Notification {
	return LowBalanceNotification{asset: asset, fiatAmount: fiatAmount, balanceAmount: balanceAmount}
}

type LowBalanceNotification struct {
	asset         model.Asset
	fiatAmount    float64
	balanceAmount float64
}

func (n LowBalanceNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: low %s balance", n.asset.Name)
}

func (n LowBalanceNotification) Body() string {
	return fmt.Sprintf(`Your balance is running low - purchases may start failing soon.

You have %s in your account, and the next order amount is %s.

Top up your account balance ASAP.`,
		util.FormatAsset(n.asset, n.balanceAmount),
		util.FormatAsset(n.asset, n.fiatAmount))
}
