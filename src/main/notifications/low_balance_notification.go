package notifications

import (
	"fmt"
)

func NewLowBalanceNotification(currency string, fiatAmount float64, balanceAmount float64) Notification {
	return LowBalanceNotification{currency: currency, fiatAmount: fiatAmount, balanceAmount: balanceAmount}
}

type LowBalanceNotification struct {
	currency      string
	fiatAmount    float64
	balanceAmount float64
}

func (n LowBalanceNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: low balance for %s", n.currency)
}

func (n LowBalanceNotification) Body() string {
	return fmt.Sprintf(`Your balance is running low - purchases may start failing soon.

You have %f %s in your account, and the next order amount is %f.

Top up your account balance ASAP.`,
		n.balanceAmount,
		n.currency,
		n.fiatAmount,
	)

}
