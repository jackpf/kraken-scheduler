package model

func NewOrder(pair string, price float64, fiatAmount float64) Order {
	return Order{Pair: pair, Price: price, FiatAmount: fiatAmount}
}

type Order struct {
	Pair       string
	Price      float64
	FiatAmount float64
}

func (o Order) Amount() float64 {
	return o.FiatAmount / o.Price
}
