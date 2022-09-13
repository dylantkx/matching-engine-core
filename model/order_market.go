package model

import "github.com/shopspring/decimal"

type OrderMarket struct {
	ID    int64           `json:"id"`
	Units decimal.Decimal `json:"units"`
	Side  OrderSide       `json:"side"`
}
