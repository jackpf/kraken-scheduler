package model

import (
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
)

type BalanceRequest struct {
	Pair   model.Pair
	Amount float64
}

type BalanceData struct {
	Asset              model.Asset
	NextPurchaseAmount float64
	Balance            float64
}
