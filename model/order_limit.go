package model

import "github.com/shopspring/decimal"

type OrderLimit struct {
	ID    int64           `json:"id"`
	Units decimal.Decimal `json:"units"`
	Price decimal.Decimal `json:"price"`
	Side  OrderSide       `json:"side"`
}
