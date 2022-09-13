package orderbook_test

import (
	"testing"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/dylantkx/matching-engine-core/orderbook"
	"github.com/shopspring/decimal"
)

func TestGetFullBookSnapshot(t *testing.T) {
	b := orderbook.NewBook()

	for i := 1; i <= 10; i++ {
		b.AddBuyOrder(model.Order{
			ID:    int64(i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Buy,
		})
	}
	for i := 11; i <= 20; i++ {
		b.AddSellOrder(model.Order{
			ID:    int64(i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Sell,
		})
	}

	sn := b.GetFullSnapshot()
	if len(sn.Buys) != 10 {
		t.Fatalf("expect to have 10 buys, but got %d", len(sn.Buys))
	}
	if len(sn.Sells) != 10 {
		t.Fatalf("expect to have 10 sells, but got %d", len(sn.Sells))
	}
}

func TestGetBookSnapshotWithDepth(t *testing.T) {
	b := orderbook.NewBook()

	for i := 1; i <= 10; i++ {
		b.AddBuyOrder(model.Order{
			ID:    int64(i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Buy,
		})
	}
	for i := 11; i <= 20; i++ {
		b.AddSellOrder(model.Order{
			ID:    int64(i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Sell,
		})
	}

	sn := b.GetSnapshotWithDepth(2)
	if len(sn.Buys) != 2 {
		t.Fatalf("expect to have 2 buys, but got %d", len(sn.Buys))
	}
	if len(sn.Sells) != 2 {
		t.Fatalf("expect to have 2 sells, but got %d", len(sn.Sells))
	}
}

func TestGetBookSnapshotWithDepth2(t *testing.T) {
	b := orderbook.NewBook()

	for i := 1; i <= 10; i++ {
		b.AddBuyOrder(model.Order{
			ID:    int64(i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Buy,
		})
	}
	for i := 11; i <= 20; i++ {
		b.AddSellOrder(model.Order{
			ID:    int64(i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Sell,
		})
	}

	sn := b.GetSnapshotWithDepth(11)
	if len(sn.Buys) != 10 {
		t.Fatalf("expect to have 10 buys, but got %d", len(sn.Buys))
	}
	if len(sn.Sells) != 10 {
		t.Fatalf("expect to have 10 sells, but got %d", len(sn.Sells))
	}
}
