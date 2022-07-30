package model

func NewOrder(pair string, price float32, fiatAmount float32) Order {
	return Order{Pair: pair, Price: price, FiatAmount: fiatAmount}
}

type Order struct {
	Pair       string
	Price      float32
	FiatAmount float32
}

func (o Order) Amount() float32 {
	return o.FiatAmount / o.Price
}
