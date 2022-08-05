package model

import "fmt"

type BalanceRequest struct {
	Pair   string
	Amount float64
}

func (r BalanceRequest) Currency() string {
	return fmt.Sprintf("%s", r.Pair[len(r.Pair)-4:])
}

type BalanceData struct {
	Currency           string
	NextPurchaseAmount float64
	Balance            float64
}
