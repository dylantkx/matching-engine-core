package model

import "github.com/shopspring/decimal"

type Order struct {
	ID    string
	Units decimal.Decimal
	Price decimal.Decimal
	Side  OrderSide
}

func (o *Order) GetVolume() decimal.Decimal {
	return o.Units.Mul(o.Price)
}

func (o *Order) Clone() Order {
	return Order{
		ID:    o.ID,
		Units: o.Units.Copy(),
		Price: o.Price.Copy(),
		Side:  o.Side,
	}
}
