package model

import "github.com/jackpf/kraken-scheduler/src/main/config/model"

func NewOrder(pair model.Pair, price float64, fiatAmount float64) Order {
	return Order{Pair: pair, Price: price, FiatAmount: fiatAmount}
}

type Order struct {
	Pair       model.Pair
	Price      float64
	FiatAmount float64
}

func (o Order) Amount() float64 {
	return o.FiatAmount / o.Price
}
