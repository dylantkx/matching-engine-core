package matchingenginecore_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	me "github.com/dylantkx/matching-engine-core"
	"github.com/dylantkx/matching-engine-core/model"
	"github.com/shopspring/decimal"
)

func TestProcessLimitOneBuyOrder(t *testing.T) {
	engine := me.NewMatchingEngine()

	ord := model.OrderLimit{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Buy,
	}
	r := engine.ProcessLimitOrder(&ord)
	sn := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) > 0 {
		t.Fatalf("expect no trades but got %d", len(r.Trades))
	}
	if sn == nil || len(sn.Buys) != 1 {
		t.Fatalf("expect order buy book size = 1 but got %d", len(sn.Buys))
	}
}

func TestProcessLimitOneSellOrder(t *testing.T) {
	engine := me.NewMatchingEngine()

	ord := model.OrderLimit{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Sell,
	}
	r := engine.ProcessLimitOrder(&ord)
	sn := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) > 0 {
		t.Fatalf("expect no trades but got %d", len(r.Trades))
	}
	if sn == nil || len(sn.Sells) != 1 {
		t.Fatalf("expect order book sell size = 1 but got %d", len(sn.Sells))
	}
}

func TestProcessLimitOrderWithoutTrades(t *testing.T) {
	engine := me.NewMatchingEngine()

	sellOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Sell,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Sell,
		},
	}
	for _, ord := range sellOrders {
		engine.ProcessLimitOrder(&ord)
	}
	sn1 := engine.GetOrderBookFullSnapshot()
	if sn1 == nil || len(sn1.Sells) != 2 {
		t.Fatalf("expect order book sell size = 2 but got %d", len(sn1.Sells))
	}

	buyOrder := model.OrderLimit{
		ID:    "3",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(50),
		Side:  model.OrderSide_Buy,
	}
	r := engine.ProcessLimitOrder(&buyOrder)
	sn2 := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) != 0 {
		t.Fatalf("expect 1 trade but got %d", len(r.Trades))
	}
	if sn2 == nil || len(sn2.Buys) != 1 {
		t.Fatalf("expect order buy book size = 1 but got %d", len(sn2.Buys))
	}
}

func TestLimitBuyOrderProducingTrades(t *testing.T) {
	engine := me.NewMatchingEngine()

	sellOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Sell,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Sell,
		},
	}
	for _, ord := range sellOrders {
		engine.ProcessLimitOrder(&ord)
	}
	sn1 := engine.GetOrderBookFullSnapshot()
	if sn1 == nil || len(sn1.Sells) != 2 {
		t.Fatalf("expect order book sell size = 2 but got %d", len(sn1.Sells))
	}

	buyOrder := model.OrderLimit{
		ID:    "3",
		Units: decimal.NewFromFloat(1.5),
		Price: decimal.NewFromFloat(200),
		Side:  model.OrderSide_Buy,
	}
	r := engine.ProcessLimitOrder(&buyOrder)
	sn2 := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) != 2 {
		t.Fatalf("expect 2 trade but got %d", len(r.Trades))
	}
	tr := r.Trades[0]
	if tr.BuyOrderID != buyOrder.ID || tr.IsBuyerMaker || tr.Price.GreaterThan(buyOrder.Price) || !tr.Units.Equal(sellOrders[0].Units) {
		t.Fatalf("wrong trade output: %+v", tr)
	}

	if sn2 == nil || len(sn2.Sells) != 1 {
		t.Fatalf("expect order sell book size = 1 but got %d", len(sn2.Sells))
	}
}

func TestLimitSellOrderProducingTrades(t *testing.T) {
	engine := me.NewMatchingEngine()

	buyOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Buy,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Buy,
		},
	}
	for _, ord := range buyOrders {
		engine.ProcessLimitOrder(&ord)
	}
	sn1 := engine.GetOrderBookFullSnapshot()
	if sn1 == nil || len(sn1.Buys) != 2 {
		t.Fatalf("expect order book buy size = 2 but got %d", len(sn1.Buys))
	}

	sellOrder := model.OrderLimit{
		ID:    "3",
		Units: decimal.NewFromFloat(2),
		Price: decimal.NewFromFloat(200),
		Side:  model.OrderSide_Sell,
	}
	r := engine.ProcessLimitOrder(&sellOrder)
	sn2 := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) != 1 {
		t.Fatalf("expect 1 trade but got %d", len(r.Trades))
	}
	tr := r.Trades[0]
	if tr.SellOrderID != sellOrder.ID || !tr.IsBuyerMaker || tr.Price.LessThan(sellOrder.Price) || !tr.Units.Equal(decimal.NewFromFloat(1)) {
		t.Fatalf("wrong trade output: %+v", tr)
	}

	if sn2 == nil || len(sn2.Sells) != 1 {
		t.Fatalf("expect order sell book size = 1 but got %d", len(sn2.Sells))
	}
	if sn2 == nil || len(sn2.Buys) != 1 {
		t.Fatalf("expect order buy book size = 1 but got %d", len(sn2.Buys))
	}
}

func TestProcessMarketOneBuyOrder(t *testing.T) {
	engine := me.NewMatchingEngine()

	ord := model.OrderMarket{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Side:  model.OrderSide_Buy,
	}
	r := engine.ProcessMarketOrder(&ord)

	if len(r.Trades) > 0 {
		t.Fatalf("expect no trades but got %d", len(r.Trades))
	}
	if len(r.Cancellations) != 1 {
		t.Fatalf("expect 1 cancellation but got %d", len(r.Cancellations))
	}
	if !r.Cancellations[0].Units.Equal(ord.Units) {
		t.Fatalf("expect to cancel whole order, but got %s", r.Cancellations[0].Units)
	}
}

func TestProcessMarketOneSellOrder(t *testing.T) {
	engine := me.NewMatchingEngine()

	ord := model.OrderMarket{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Side:  model.OrderSide_Sell,
	}
	r := engine.ProcessMarketOrder(&ord)

	if len(r.Trades) > 0 {
		t.Fatalf("expect no trades but got %d", len(r.Trades))
	}
	if len(r.Cancellations) != 1 {
		t.Fatalf("expect 1 cancellation but got %d", len(r.Cancellations))
	}
	if !r.Cancellations[0].Units.Equal(ord.Units) {
		t.Fatalf("expect to cancel whole order, but got %s", r.Cancellations[0].Units)
	}
}

func TestMarketBuyOrderProducingTrades(t *testing.T) {
	engine := me.NewMatchingEngine()

	sellOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Sell,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Sell,
		},
	}
	for _, ord := range sellOrders {
		engine.ProcessLimitOrder(&ord)
	}

	buyOrder := model.OrderMarket{
		ID:    "3",
		Units: decimal.NewFromFloat(1.5),
		Side:  model.OrderSide_Buy,
	}
	r := engine.ProcessMarketOrder(&buyOrder)
	sn2 := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) != 2 {
		t.Fatalf("expect 2 trades but got %d", len(r.Trades))
	}
	if len(r.Cancellations) != 0 {
		t.Fatalf("expect 0 cancels but got %d", len(r.Cancellations))
	}
	tr1, tr2 := r.Trades[0], r.Trades[1]
	if tr1.BuyOrderID != buyOrder.ID || tr1.IsBuyerMaker || !tr1.Price.Equal(sellOrders[0].Price) || !tr1.Units.Equal(sellOrders[0].Units) {
		t.Fatalf("wrong trade output: %+v", tr1)
	}
	if tr2.BuyOrderID != buyOrder.ID || tr2.IsBuyerMaker || !tr2.Price.Equal(sellOrders[1].Price) || !tr2.Units.Equal(buyOrder.Units.Sub(sellOrders[0].Units)) {
		t.Fatalf("wrong trade output: %+v", tr2)
	}

	if sn2 == nil || len(sn2.Sells) != 1 {
		t.Fatalf("expect order sell book size = 1 but got %d", len(sn2.Sells))
	}
}

func TestMarketSellOrderProducingTrades(t *testing.T) {
	engine := me.NewMatchingEngine()

	buyOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Buy,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Buy,
		},
	}
	for _, ord := range buyOrders {
		engine.ProcessLimitOrder(&ord)
	}

	sellOrder := model.OrderMarket{
		ID:    "3",
		Units: decimal.NewFromFloat(1.5),
		Side:  model.OrderSide_Sell,
	}
	r := engine.ProcessMarketOrder(&sellOrder)
	sn2 := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) != 2 {
		t.Fatalf("expect 2 trade but got %d", len(r.Trades))
	}
	if len(r.Cancellations) != 0 {
		t.Fatalf("expect 0 cancels but got %d", len(r.Cancellations))
	}
	tr1, tr2 := r.Trades[0], r.Trades[1]
	if tr1.SellOrderID != sellOrder.ID || !tr1.IsBuyerMaker || !tr1.Price.Equal(buyOrders[0].Price) || !tr1.Units.Equal(buyOrders[0].Units) {
		t.Fatalf("wrong trade output: %+v", tr1)
	}
	if tr2.SellOrderID != sellOrder.ID || !tr2.IsBuyerMaker || !tr2.Price.Equal(buyOrders[1].Price) || !tr2.Units.Equal(sellOrder.Units.Sub(buyOrders[0].Units)) {
		t.Fatalf("wrong trade output: %+v", tr2)
	}

	if sn2 == nil || len(sn2.Buys) != 1 {
		t.Fatalf("expect order buy book size = 1 but got %d", len(sn2.Buys))
	}
}

func TestMarketBuyOrderProducingTradesWithCancels(t *testing.T) {
	engine := me.NewMatchingEngine()

	sellOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Sell,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Sell,
		},
	}
	for _, ord := range sellOrders {
		engine.ProcessLimitOrder(&ord)
	}

	buyOrder := model.OrderMarket{
		ID:    "3",
		Units: decimal.NewFromFloat(2.5),
		Side:  model.OrderSide_Buy,
	}
	r := engine.ProcessMarketOrder(&buyOrder)
	sn2 := engine.GetOrderBookFullSnapshot()

	if len(r.Trades) != 2 {
		t.Fatalf("expect 2 trade but got %d", len(r.Trades))
	}
	if len(r.Cancellations) != 1 {
		t.Fatalf("expect 1 cancels but got %d", len(r.Cancellations))
	}
	tr1, tr2 := r.Trades[0], r.Trades[1]
	if tr1.BuyOrderID != buyOrder.ID || tr1.IsBuyerMaker || !tr1.Price.Equal(sellOrders[0].Price) || !tr1.Units.Equal(sellOrders[0].Units) {
		t.Fatalf("wrong trade output: %+v", tr1)
	}
	if tr2.BuyOrderID != buyOrder.ID || tr2.IsBuyerMaker || !tr2.Price.Equal(sellOrders[1].Price) || !tr2.Units.Equal(sellOrders[0].Units) {
		t.Fatalf("wrong trade output: %+v", tr2)
	}
	if !r.Cancellations[0].Units.Equal(buyOrder.Units.Sub(sellOrders[0].Units.Add(sellOrders[1].Units))) {
		t.Fatalf("wrong cancelled units, got %s", r.Cancellations[0].Units)
	}

	if sn2 == nil || len(sn2.Sells) != 0 {
		fmt.Printf("%+v\n", sn2.Sells)
		t.Fatalf("expect order sell book size = 0 but got %d", len(sn2.Sells))
	}
}

func TestNearlyConcurrentMarketBuys(t *testing.T) {
	engine := me.NewMatchingEngine()

	sellOrders := []model.OrderLimit{
		{
			ID:    "1",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(100),
			Side:  model.OrderSide_Sell,
		},
		{
			ID:    "2",
			Units: decimal.NewFromFloat(1),
			Price: decimal.NewFromFloat(200),
			Side:  model.OrderSide_Sell,
		},
	}
	for _, ord := range sellOrders {
		engine.ProcessLimitOrder(&ord)
	}

	buyOrder1 := model.OrderMarket{
		ID:    "3",
		Units: decimal.NewFromFloat(1),
		Side:  model.OrderSide_Buy,
	}
	buyOrder2 := model.OrderMarket{
		ID:    "4",
		Units: decimal.NewFromFloat(1),
		Side:  model.OrderSide_Buy,
	}

	var trades_1, trades_2 []model.Trade

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		r := engine.ProcessMarketOrder(&buyOrder1)
		trades_1 = r.Trades
	}()
	go func() {
		time.Sleep(time.Millisecond * 1)
		defer wg.Done()
		r := engine.ProcessMarketOrder(&buyOrder2)
		trades_2 = r.Trades
	}()

	wg.Wait()

	if len(trades_1) != 1 {
		t.Fatalf("expect 1 trade but got %d", len(trades_1))
	}

	if !trades_1[0].Price.Equal(sellOrders[0].Price) {
		t.Fatalf("expect trade 1 to be filled with price of %s, but got %s", sellOrders[0].Price, trades_1[0].Price)
	}
	if !trades_2[0].Price.Equal(sellOrders[1].Price) {
		t.Fatalf("expect trade 2 to be filled with price of %s, but got %s", sellOrders[1].Price, trades_2[0].Price)
	}
}

func TestCancelOrder(t *testing.T) {
	engine := me.NewMatchingEngine()

	ord := model.OrderLimit{
		ID:    "1",
		Units: decimal.NewFromFloat(1),
		Price: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Buy,
	}
	engine.ProcessLimitOrder(&ord)

	cancels, err := engine.CancelOrder(model.Order{
		ID:    ord.ID,
		Units: ord.Units.Copy(),
		Price: ord.Price.Copy(),
		Side:  ord.Side,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	sn := engine.GetOrderBookFullSnapshot()

	if len(cancels) != 1 {
		t.Fatalf("expect 1 cancel but got %d", len(cancels))
	}
	if cancels[0].OrderID != ord.ID || !cancels[0].Units.Equal(ord.Units) {
		t.Fatalf("wrong cancel output: %+v", cancels)
	}
	if sn == nil || len(sn.Buys) != 0 {
		t.Fatalf("expect order buy book size = 0 but got %d", len(sn.Buys))
	}
}

func BenchmarkProcessLimitOrders(b *testing.B) {
	b.StopTimer()
	engine := me.NewMatchingEngine()
	orders := make([]*model.OrderLimit, 0, b.N)
	for i := 0; i < b.N; i++ {
		order := &model.OrderLimit{
			ID:    fmt.Sprintf("%d", i+1),
			Units: decimal.NewFromFloat((rand.Float64() + float64(rand.Intn(2))) * 10),
			Price: decimal.NewFromFloat((rand.Float64() + float64(rand.Intn(2))) * 100),
			Side:  model.OrderSide_Buy,
		}
		orders = append(orders, order)
	}
	b.StartTimer()
	for _, order := range orders {
		engine.ProcessLimitOrder(order)
	}
}

func BenchmarkProcessLimitOrdersAsync(b *testing.B) {
	b.StopTimer()
	engine := me.NewMatchingEngine()
	orders := make([]*model.OrderLimit, 0, b.N)
	for i := 0; i < b.N; i++ {
		order := &model.OrderLimit{
			ID:    fmt.Sprintf("%d", i+1),
			Units: decimal.NewFromFloat((rand.Float64() + float64(rand.Intn(2))) * 10),
			Price: decimal.NewFromFloat((rand.Float64() + float64(rand.Intn(2))) * 100),
			Side:  model.OrderSide_Buy,
		}
		orders = append(orders, order)
	}
	b.StartTimer()
	wg := sync.WaitGroup{}
	for _, order := range orders {
		wg.Add(1)
		go func(order *model.OrderLimit) {
			defer wg.Done()
			engine.ProcessLimitOrder(order)
		}(order)
	}
	wg.Wait()
}

func BenchmarkProcessOneMarketOrderProducingTrades(b *testing.B) {
	b.StopTimer()
	engine := me.NewMatchingEngine()
	var n int = 1e5
	orders := make([]*model.OrderLimit, 0, n)
	for i := 0; i < n; i++ {
		order := &model.OrderLimit{
			ID:    fmt.Sprintf("%d", i+1),
			Units: decimal.NewFromFloat((rand.Float64() + float64(rand.Intn(2))) * 10),
			Price: decimal.NewFromFloat((rand.Float64() + float64(rand.Intn(2))) * 100),
			Side:  model.OrderSide_Buy,
		}
		orders = append(orders, order)
	}
	for _, order := range orders {
		engine.ProcessLimitOrder(order)
	}

	order := &model.OrderMarket{
		ID:    fmt.Sprintf("%d", n+1),
		Units: decimal.NewFromFloat(100),
		Side:  model.OrderSide_Sell,
	}

	b.StartTimer()
	engine.ProcessMarketOrder(order)
}
