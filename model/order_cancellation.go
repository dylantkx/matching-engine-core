package model

import "github.com/shopspring/decimal"

type OrderCancellation struct {
	OrderID string          `json:"orderId"`
	Units   decimal.Decimal `json:"units"`
}
