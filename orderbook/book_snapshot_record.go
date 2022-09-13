package orderbook

import (
	"github.com/shopspring/decimal"
)

type bookSnapshotRecord struct {
	Price decimal.Decimal `json:"price"`
	Size  decimal.Decimal `json:"size"`
}

func NewBookSnapshotRecord(price, size decimal.Decimal) *bookSnapshotRecord {
	return &bookSnapshotRecord{
		Price: price,
		Size:  size,
	}
}
