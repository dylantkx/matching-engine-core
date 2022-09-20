package model

import "github.com/shopspring/decimal"

type OrderMarket struct {
	ID    string          `json:"id"`
	Units decimal.Decimal `json:"units"`
	Side  OrderSide       `json:"side"`
}
