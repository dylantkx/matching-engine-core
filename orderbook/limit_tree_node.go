package orderbook

import "github.com/shopspring/decimal"

type limitTreeNode struct {
	Price    decimal.Decimal
	LimitRef *bookLimit
}
