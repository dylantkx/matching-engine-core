package model

import "github.com/shopspring/decimal"

type OrderCancellation struct {
	OrderID int64           `json:"orderId"`
	Units   decimal.Decimal `json:"units"`
}
