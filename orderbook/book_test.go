package orderbook_test

import (
	"fmt"
	"testing"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/dylantkx/matching-engine-core/orderbook"
	"github.com/shopspring/decimal"
)

func TestGetFullBookSnapshot(t *testing.T) {
	b := orderbook.NewBook()

	for i := 1; i <= 10; i++ {
		b.AddBuyOrder(model.Order{
			ID:    fmt.Sprintf("%d", i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Buy,
		})
	}
	for i := 11; i <= 20; i++ {
		b.AddSellOrder(model.Order{
			ID:    fmt.Sprintf("%d", i),
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
			ID:    fmt.Sprintf("%d", i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Buy,
		})
	}
	for i := 11; i <= 20; i++ {
		b.AddSellOrder(model.Order{
			ID:    fmt.Sprintf("%d", i),
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
			ID:    fmt.Sprintf("%d", i),
			Units: decimal.NewFromFloat(float64(i)),
			Price: decimal.NewFromFloat(float64(i) * 100),
			Side:  model.OrderSide_Buy,
		})
	}
	for i := 11; i <= 20; i++ {
		b.AddSellOrder(model.Order{
			ID:    fmt.Sprintf("%d", i),
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

func TestGetTotalBuyUnitsFromPrice(t *testing.T) {
	b := orderbook.NewBook()

	b.AddBuyOrder(model.Order{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(50),
		Side:  model.OrderSide_Buy,
	})
	b.AddBuyOrder(model.Order{
		ID:    "2",
		Units: decimal.NewFromFloat(2),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Buy,
	})
	b.AddBuyOrder(model.Order{
		ID:    "3",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Buy,
	})

	units := b.GetTotalBuyUnitsFromPrice(decimal.NewFromFloat(100))
	if !units.Equal(decimal.NewFromFloat(3)) {
		t.Fatalf("expect to units = 3, but got %s", units)
	}

	units = b.GetTotalBuyUnitsFromPrice(decimal.NewFromFloat(70))
	if !units.Equal(decimal.NewFromFloat(3)) {
		t.Fatalf("expect to units = 3, but got %s", units)
	}
	units = b.GetTotalBuyUnitsFromPrice(decimal.NewFromFloat(50))
	if !units.Equal(decimal.NewFromFloat(4)) {
		t.Fatalf("expect to units = 4, but got %s", units)
	}
}

func TestGetTotalSellUnitsToPrice(t *testing.T) {
	b := orderbook.NewBook()

	b.AddSellOrder(model.Order{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(50),
		Side:  model.OrderSide_Sell,
	})
	b.AddSellOrder(model.Order{
		ID:    "2",
		Units: decimal.NewFromFloat(2),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Sell,
	})
	b.AddSellOrder(model.Order{
		ID:    "3",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Sell,
	})

	units := b.GetTotalSellUnitsToPrice(decimal.NewFromFloat(100))
	if !units.Equal(decimal.NewFromFloat(4)) {
		t.Fatalf("expect to units = 4, but got %s", units)
	}

	units = b.GetTotalSellUnitsToPrice(decimal.NewFromFloat(70))
	if !units.Equal(decimal.NewFromFloat(1)) {
		t.Fatalf("expect to units = 1, but got %s", units)
	}
	units = b.GetTotalSellUnitsToPrice(decimal.NewFromFloat(50))
	if !units.Equal(decimal.NewFromFloat(1)) {
		t.Fatalf("expect to units = 1, but got %s", units)
	}
}
