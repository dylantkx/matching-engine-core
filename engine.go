package matchingenginecore

import (
	"time"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/dylantkx/matching-engine-core/orderbook"
	"github.com/shopspring/decimal"
)

type MatchingEngine struct {
	book orderbook.Book
}

func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		book: orderbook.NewBook(),
	}
}

func (me *MatchingEngine) GetHighestBuyPrice() decimal.Decimal {
	best := me.book.GetHighestBuy()
	if best == nil {
		return decimal.Zero
	}
	return best.Price
}

func (me *MatchingEngine) GetLowestSellPrice() decimal.Decimal {
	best := me.book.GetLowestSell()
	if best == nil {
		return decimal.Zero
	}
	return best.Price
}

func (me *MatchingEngine) GetOrderBookFullSnapshot() *orderbook.BookSnapshot {
	return me.book.GetFullSnapshot()
}

func (me *MatchingEngine) GetOrderBookSnapshotWithDepth(depth int) *orderbook.BookSnapshot {
	return me.book.GetSnapshotWithDepth(depth)
}

func (me *MatchingEngine) GetTotalBuyUnitsFromPrice(price decimal.Decimal) decimal.Decimal {
	return me.book.GetTotalBuyUnitsFromPrice(price)
}

func (me *MatchingEngine) GetTotalSellUnitsToPrice(price decimal.Decimal) decimal.Decimal {
	return me.book.GetTotalSellUnitsToPrice(price)
}

func (me *MatchingEngine) ProcessLimitOrder(order *model.OrderLimit) model.MatchResult {
	if order.Side == model.OrderSide_Buy {
		return me.processLimitBuyOrder(order)
	}
	return me.processLimitSellOrder(order)
}

func (me *MatchingEngine) ProcessMarketOrder(order *model.OrderMarket) model.MatchResult {
	if order.Side == model.OrderSide_Buy {
		return me.processMarketBuyOrder(order)
	}
	return me.processMarketSellOrder(order)
}

func (me *MatchingEngine) CancelOrder(order model.Order) ([]model.OrderCancellation, error) {
	return me.book.CancelOrder(order)
}

func (me *MatchingEngine) processLimitBuyOrder(order *model.OrderLimit) (r model.MatchResult) {
	if me.book.GetLowestSell() == nil || me.book.GetLowestSell().Price.GreaterThan(order.Price) {
		me.book.AddBuyOrder(model.Order{
			ID:    order.ID,
			Units: order.Units,
			Price: order.Price,
			Side:  model.OrderSide_Buy,
		})
		return
	}

	now := time.Now()
	remainingUnits := order.Units.Copy()

	matchedOrders := me.book.ClearSellSideByUnitsAndPrice(order.Units.Copy(), order.Price.Copy())
	for _, o := range matchedOrders {
		r.Trades = append(r.Trades, model.Trade{
			BuyOrderID:   order.ID,
			SellOrderID:  o.ID,
			Units:        o.Units,
			Price:        o.Price,
			IsBuyerMaker: false,
			EventTime:    model.Timestamp{Time: now},
		})
		remainingUnits = remainingUnits.Sub(o.Units)
	}

	// push the order into buy book if any remaining
	if remainingUnits.IsPositive() {
		me.book.AddBuyOrder(model.Order{
			ID:    order.ID,
			Units: remainingUnits,
			Price: order.Price,
			Side:  model.OrderSide_Buy,
		})
	}
	return
}

func (me *MatchingEngine) processLimitSellOrder(order *model.OrderLimit) (r model.MatchResult) {
	if me.book.GetHighestBuy() == nil || me.book.GetHighestBuy().Price.LessThan(order.Price) {
		me.book.AddSellOrder(model.Order{
			ID:    order.ID,
			Units: order.Units,
			Price: order.Price,
			Side:  model.OrderSide_Sell,
		})
		return
	}

	now := time.Now()
	remainingUnits := order.Units.Copy()

	matchedOrders := me.book.ClearBuySideByUnitsAndPrice(order.Units.Copy(), order.Price.Copy())
	for _, o := range matchedOrders {
		r.Trades = append(r.Trades, model.Trade{
			BuyOrderID:   o.ID,
			SellOrderID:  order.ID,
			Units:        o.Units,
			Price:        o.Price,
			IsBuyerMaker: true,
			EventTime:    model.Timestamp{Time: now},
		})
		remainingUnits = remainingUnits.Sub(o.Units)
	}

	// push the order into sell book if any remaining
	if remainingUnits.IsPositive() {
		me.book.AddSellOrder(model.Order{
			ID:    order.ID,
			Units: remainingUnits,
			Price: order.Price,
			Side:  model.OrderSide_Sell,
		})
	}
	return
}

func (me *MatchingEngine) processMarketBuyOrder(order *model.OrderMarket) (r model.MatchResult) {
	if me.book.GetLowestSell() == nil {
		r.Cancellations = append(r.Cancellations, model.OrderCancellation{
			OrderID: order.ID,
			Units:   order.Units,
		})
		return
	}

	now := time.Now()
	remainingUnits := order.Units.Copy()

	matchedOrders := me.book.ClearSellSideByUnits(order.Units.Copy())
	for _, o := range matchedOrders {
		r.Trades = append(r.Trades, model.Trade{
			BuyOrderID:   order.ID,
			SellOrderID:  o.ID,
			Units:        o.Units,
			Price:        o.Price,
			IsBuyerMaker: false,
			EventTime:    model.Timestamp{Time: now},
		})
		remainingUnits = remainingUnits.Sub(o.Units)
	}

	if remainingUnits.IsPositive() {
		r.Cancellations = append(r.Cancellations, model.OrderCancellation{
			OrderID: order.ID,
			Units:   remainingUnits,
		})
	}
	return
}

func (me *MatchingEngine) processMarketSellOrder(order *model.OrderMarket) (r model.MatchResult) {
	if me.book.GetHighestBuy() == nil {
		r.Cancellations = append(r.Cancellations, model.OrderCancellation{
			OrderID: order.ID,
			Units:   order.Units,
		})
		return
	}

	now := time.Now()
	remainingUnits := order.Units.Copy()

	matchedOrders := me.book.ClearBuySideByUnits(order.Units.Copy())
	for _, o := range matchedOrders {
		r.Trades = append(r.Trades, model.Trade{
			BuyOrderID:   o.ID,
			SellOrderID:  order.ID,
			Units:        o.Units,
			Price:        o.Price,
			IsBuyerMaker: true,
			EventTime:    model.Timestamp{Time: now},
		})
		remainingUnits = remainingUnits.Sub(o.Units)
	}

	if remainingUnits.IsPositive() {
		r.Cancellations = append(r.Cancellations, model.OrderCancellation{
			OrderID: order.ID,
			Units:   remainingUnits,
		})
	}
	return
}
