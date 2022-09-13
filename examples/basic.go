package main

import (
	"fmt"

	me "github.com/dylantkx/matching-engine-core"
	"github.com/dylantkx/matching-engine-core/model"
	"github.com/shopspring/decimal"
)

func main() {
	engine := me.NewMatchingEngine()

	sellOrders := []model.OrderLimit{
		{
			ID:    1,
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Sell,
		},
		{
			ID:    2,
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Sell,
		},
	}
	for _, ord := range sellOrders {
		engine.ProcessLimitOrder(&ord)
	}

	sn1 := engine.GetOrderBookFullSnapshot()
	fmt.Printf("Order book snapshot 1: %+v\n", sn1)

	buyOrder := model.OrderMarket{
		ID:    3,
		Units: decimal.NewFromFloat(1.5),
		Side:  model.OrderSide_Buy,
	}
	trades, cancels := engine.ProcessMarketOrder(&buyOrder)
	fmt.Printf("Trades: %+v\n", trades)
	fmt.Printf("Cancels: %+v\n", cancels)

	sn2 := engine.GetOrderBookFullSnapshot()
	fmt.Printf("Order book snapshot 2: %+v\n", sn2)
}
