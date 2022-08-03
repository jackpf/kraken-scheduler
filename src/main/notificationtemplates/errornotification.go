package notificationtemplates

import (
	"fmt"
)

func NewErrorNotification(pair string, amount float64, orderPrice float64, err error) ErrorNotification {
	return ErrorNotification{pair: pair, amount: amount, orderPrice: orderPrice, err: err}
}

type ErrorNotification struct {
	pair       string
	amount     float64
	orderPrice float64
	err        error
}

func (n ErrorNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: %s order failed", n.pair)
}

func (n ErrorNotification) Body() string {
	return fmt.Sprintf(`An error occured attempting to order %f %s, at a cost of %f.

Error: %s`,
		n.amount,
		n.pair,
		n.orderPrice,
		n.err.Error())

}
