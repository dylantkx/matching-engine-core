package orderbook_test

import (
	"testing"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/dylantkx/matching-engine-core/orderbook"
	"github.com/shopspring/decimal"
)

func TestBookLimitInsertOrUpdateOrder(t *testing.T) {
	bl := orderbook.NewBookLimit()

	order := &model.Order{
		ID:    1,
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Buy,
	}
	isUpdated := bl.InsertOrUpdateOrder(order)

	if isUpdated {
		t.Fatalf("expect to be insertion, not update")
	}
	if !bl.Size.Equal(order.Units) {
		t.Fatalf("expect size to be %s, got %s", order.Units, bl.Size)
	}
	if !bl.Price.Equal(order.Price) {
		t.Fatalf("expect price to be %s, got %s", order.Price, bl.Price)
	}
	if !bl.Volume.Equal(order.GetVolume()) {
		t.Fatalf("expect volume to be %s, got %s", order.GetVolume(), bl.Volume)
	}
}

func TestBookLimitRemoveOrder(t *testing.T) {
	bl := orderbook.NewBookLimit()

	order := &model.Order{
		ID:    1,
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Buy,
	}
	bl.InsertOrUpdateOrder(order)

	bl.RemoveOrder(order)

	if !bl.IsEmpty() {
		t.Fatalf("expect book limit to be empty")
	}
	if !bl.Volume.IsZero() {
		t.Fatalf("expect volume to be zero, but got %+v", bl.Volume)
	}
	if !bl.Size.IsZero() {
		t.Fatalf("expect size to be zero, but got %+v", bl.Size)
	}
}

func TestBookLimitRemoveOrder2(t *testing.T) {
	bl := orderbook.NewBookLimit()

	orders := make([]*model.Order, 0, 10)
	for i := 0; i < 10; i++ {
		order := &model.Order{
			ID:    int64(i + 1),
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Buy,
		}
		orders = append(orders, order)
		bl.InsertOrUpdateOrder(order)
	}

	for idx, ord := range orders {
		if idx >= 5 {
			break
		}
		bl.RemoveOrder(ord)
	}

	if bl.CountOrders() != 5 {
		t.Fatalf("expect book limit to have 5 orders, but got %d", bl.CountOrders())
	}
}
