package model

import (
	"github.com/shopspring/decimal"
)

type Trade struct {
	BuyOrderID   string          `json:"buyOrderId"`
	SellOrderID  string          `json:"sellOrderId"`
	Units        decimal.Decimal `json:"units"`
	Price        decimal.Decimal `json:"price"`
	IsBuyerMaker bool            `json:"isBuyerMaker"`
	EventTime    Timestamp       `json:"eventTime"`
}
