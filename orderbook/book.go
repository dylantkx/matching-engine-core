package orderbook

import (
	"errors"
	"sync"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/google/btree"
	"github.com/shopspring/decimal"
)

const (
	treeDegree       int = 2
	maxOrderPerLimit int = 1e5
)

type Book interface {
	AddBuyOrder(order model.Order)
	AddSellOrder(order model.Order)
	CancelOrder(order model.Order) error
	ClearBuySideByUnitsAndPrice(units decimal.Decimal, price decimal.Decimal) (clearedOrders []*model.Order)
	ClearSellSideByUnitsAndPrice(units decimal.Decimal, price decimal.Decimal) (clearedOrders []*model.Order)
	ClearBuySideByUnits(units decimal.Decimal) (clearedOrders []*model.Order)
	ClearSellSideByUnits(units decimal.Decimal) (clearedOrders []*model.Order)
	GetFullSnapshot() *BookSnapshot
	GetSnapshotWithDepth(depth int) *BookSnapshot
	GetHighestBuy() *bookLimit
	GetLowestSell() *bookLimit
}

type book struct {
	buyTree     *limitTree
	buyLimitMap map[string]*bookLimit
	highestBuy  *bookLimit
	buyMu       sync.RWMutex

	sellTree     *limitTree
	sellLimitMap map[string]*bookLimit
	lowestSell   *bookLimit
	sellMu       sync.RWMutex
}

func NewBook() *book {
	return &book{
		buyTree: btree.NewWithFreeListG(treeDegree, func(a, b limitTreeNode) bool {
			return a.Price.LessThan(b.Price)
		}, btree.NewFreeListG[limitTreeNode](maxOrderPerLimit)),
		buyLimitMap: make(map[string]*bookLimit),
		sellTree: btree.NewWithFreeListG(treeDegree, func(a, b limitTreeNode) bool {
			return a.Price.LessThan(b.Price)
		}, btree.NewFreeListG[limitTreeNode](maxOrderPerLimit)),
		sellLimitMap: make(map[string]*bookLimit),
	}
}

func (b *book) GetHighestBuy() *bookLimit {
	b.buyMu.RLock()
	defer b.buyMu.RUnlock()
	return b.highestBuy
}

func (b *book) GetLowestSell() *bookLimit {
	b.sellMu.RLock()
	defer b.sellMu.RUnlock()
	return b.lowestSell
}

func (b *book) AddBuyOrder(order model.Order) {
	b.buyMu.Lock()
	defer b.buyMu.Unlock()
	bl := b.buyLimitMap[order.Price.String()]
	if bl == nil {
		bl = NewBookLimit()
		b.buyLimitMap[order.Price.String()] = bl
	}
	if isUpdate := bl.InsertOrUpdateOrder(&order); !isUpdate {
		n := limitTreeNode{
			Price:    order.Price,
			LimitRef: bl,
		}
		b.buyTree.ReplaceOrInsert(n)
		if b.highestBuy == nil || b.highestBuy.Price.LessThan(bl.Price) {
			b.highestBuy = bl
		}
	}
}

func (b *book) AddSellOrder(order model.Order) {
	b.sellMu.Lock()
	defer b.sellMu.Unlock()
	bl := b.sellLimitMap[order.Price.String()]
	if bl == nil {
		bl = NewBookLimit()
		b.sellLimitMap[order.Price.String()] = bl
	}
	if isUpdate := bl.InsertOrUpdateOrder(&order); !isUpdate {
		n := limitTreeNode{
			Price:    order.Price,
			LimitRef: bl,
		}
		b.sellTree.ReplaceOrInsert(n)
		if b.lowestSell == nil || b.lowestSell.Price.GreaterThan(bl.Price) {
			b.lowestSell = bl
		}
	}
}

func (b *book) CancelOrder(order model.Order) error {
	var t *limitTree
	var m map[string]*bookLimit
	var best *bookLimit
	var mu *sync.RWMutex
	if order.Side == model.OrderSide_Buy {
		t = b.buyTree
		m = b.buyLimitMap
		best = b.highestBuy
		mu = &b.buyMu
	} else {
		t = b.sellTree
		m = b.sellLimitMap
		best = b.lowestSell
		mu = &b.sellMu
	}
	bl := m[order.Price.String()]
	if bl == nil {
		return errors.New("order not found")
	}
	mu.Lock()
	defer mu.Unlock()
	bl.RemoveOrder(&order)
	if bl.IsEmpty() {
		delete(m, bl.Price.String())
		t.Delete(limitTreeNode{Price: bl.Price})
		if best != nil && best.Price.Equal(bl.Price) {
			best = nil
		}
	}
	return nil
}

func (b *book) ClearBuySideByUnitsAndPrice(units decimal.Decimal, price decimal.Decimal) (clearedOrders []*model.Order) {
	b.buyMu.Lock()
	defer b.buyMu.Unlock()
	clearedPrice := make([]decimal.Decimal, 0)
	b.buyTree.Descend(func(item limitTreeNode) bool {
		if item.LimitRef == nil || item.Price.LessThan(price) {
			return false
		}
		o := item.LimitRef.firstBookOrder
		for o != nil && units.IsPositive() {
			order := o.Order.Clone()
			if o.Order.Units.LessThanOrEqual(units) {
				clearedOrders = append(clearedOrders, &order)
				units = units.Sub(o.Order.Units)
				item.LimitRef.RemoveOrder(o.Order)
			} else {
				order.Units = units.Copy()
				o.Order.Units = o.Order.Units.Sub(units)
				item.LimitRef.InsertOrUpdateOrder(o.Order)
				clearedOrders = append(clearedOrders, &order)
				units = decimal.Zero
			}
			o = o.nextBookOrder
		}
		if units.IsPositive() {
			clearedPrice = append(clearedPrice, item.Price)
			return true
		}
		return false
	})
	// TODO: optimize these operations
	for _, p := range clearedPrice {
		b.buyTree.Delete(limitTreeNode{Price: p})
		delete(b.buyLimitMap, p.String())
	}
	if len(clearedPrice) > 0 {
		n, _ := b.buyTree.Max()
		b.highestBuy = n.LimitRef
	}
	return
}

func (b *book) ClearSellSideByUnitsAndPrice(units decimal.Decimal, price decimal.Decimal) (clearedOrders []*model.Order) {
	b.sellMu.Lock()
	defer b.sellMu.Unlock()
	clearedPrice := make([]decimal.Decimal, 0)
	b.sellTree.Ascend(func(item limitTreeNode) bool {
		if item.LimitRef == nil || item.Price.GreaterThan(price) {
			return false
		}
		o := item.LimitRef.firstBookOrder
		for o != nil && units.IsPositive() {
			order := o.Order.Clone()
			if o.Order.Units.LessThanOrEqual(units) {
				clearedOrders = append(clearedOrders, &order)
				units = units.Sub(o.Order.Units)
				item.LimitRef.RemoveOrder(o.Order)
			} else {
				order.Units = units.Copy()
				o.Order.Units = o.Order.Units.Sub(units)
				item.LimitRef.InsertOrUpdateOrder(o.Order)
				clearedOrders = append(clearedOrders, &order)
				units = decimal.Zero
			}
			o = o.nextBookOrder
		}
		if units.IsPositive() {
			clearedPrice = append(clearedPrice, item.Price)
			return true
		}
		return false
	})
	// TODO: optimize these operations
	for _, p := range clearedPrice {
		b.sellTree.Delete(limitTreeNode{Price: p})
		delete(b.sellLimitMap, p.String())
	}
	if len(clearedPrice) > 0 {
		n, _ := b.sellTree.Min()
		b.lowestSell = n.LimitRef
	}
	return
}

func (b *book) ClearBuySideByUnits(units decimal.Decimal) (clearedOrders []*model.Order) {
	b.buyMu.Lock()
	defer b.buyMu.Unlock()
	clearedPrice := make([]decimal.Decimal, 0)
	b.buyTree.Descend(func(item limitTreeNode) bool {
		if item.LimitRef == nil {
			return false
		}
		o := item.LimitRef.firstBookOrder
		for o != nil && units.IsPositive() {
			order := o.Order.Clone()
			if o.Order.Units.LessThanOrEqual(units) {
				clearedOrders = append(clearedOrders, &order)
				units = units.Sub(o.Order.Units)
				item.LimitRef.RemoveOrder(o.Order)
			} else {
				order.Units = units.Copy()
				o.Order.Units = o.Order.Units.Sub(units)
				item.LimitRef.InsertOrUpdateOrder(o.Order)
				clearedOrders = append(clearedOrders, &order)
				units = decimal.Zero
			}
			o = o.nextBookOrder
		}
		if units.IsPositive() {
			clearedPrice = append(clearedPrice, item.Price)
			return true
		}
		return false
	})
	// TODO: optimize these operations
	for _, p := range clearedPrice {
		b.buyTree.Delete(limitTreeNode{Price: p})
		delete(b.buyLimitMap, p.String())
	}
	if len(clearedPrice) > 0 {
		n, _ := b.buyTree.Max()
		b.highestBuy = n.LimitRef
	}
	return
}

func (b *book) ClearSellSideByUnits(units decimal.Decimal) (clearedOrders []*model.Order) {
	b.sellMu.Lock()
	defer b.sellMu.Unlock()
	clearedPrice := make([]decimal.Decimal, 0)
	b.sellTree.Ascend(func(item limitTreeNode) bool {
		if item.LimitRef == nil {
			return false
		}
		o := item.LimitRef.firstBookOrder
		for o != nil && units.IsPositive() {
			order := o.Order.Clone()
			if o.Order.Units.LessThanOrEqual(units) {
				clearedOrders = append(clearedOrders, &order)
				units = units.Sub(o.Order.Units)
				item.LimitRef.RemoveOrder(o.Order)
			} else {
				order.Units = units.Copy()
				o.Order.Units = o.Order.Units.Sub(units)
				item.LimitRef.InsertOrUpdateOrder(o.Order)
				clearedOrders = append(clearedOrders, &order)
				units = decimal.Zero
			}
			o = o.nextBookOrder
		}
		if units.IsPositive() {
			clearedPrice = append(clearedPrice, item.Price)
			return true
		}
		return false
	})
	// TODO: optimize these operations
	for _, p := range clearedPrice {
		b.sellTree.Delete(limitTreeNode{Price: p})
		delete(b.sellLimitMap, p.String())
	}
	if len(clearedPrice) > 0 {
		n, _ := b.sellTree.Min()
		b.lowestSell = n.LimitRef
	}
	return
}

func (b *book) GetFullSnapshot() *BookSnapshot {
	sn := NewBookSnapshot()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, b *book) {
		defer wg.Done()
		b.buyMu.RLock()
		defer b.buyMu.RUnlock()
		b.buyTree.Descend(func(item limitTreeNode) bool {
			if item.LimitRef == nil {
				return false
			}
			sn.Buys = append(sn.Buys, NewBookSnapshotRecord(item.LimitRef.Price, item.LimitRef.Size))
			return true
		})
	}(&wg, b)
	wg.Add(1)
	go func(wg *sync.WaitGroup, b *book) {
		defer wg.Done()
		b.sellMu.RLock()
		defer b.sellMu.RUnlock()
		b.sellTree.Ascend(func(item limitTreeNode) bool {
			if item.LimitRef == nil {
				return false
			}
			sn.Sells = append(sn.Sells, NewBookSnapshotRecord(item.LimitRef.Price, item.LimitRef.Size))
			return true
		})
	}(&wg, b)
	wg.Wait()
	return sn
}

func (b *book) GetSnapshotWithDepth(depth int) *BookSnapshot {
	sn := NewBookSnapshot()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, b *book, depth int) {
		defer wg.Done()
		b.buyMu.RLock()
		defer b.buyMu.RUnlock()
		count := 0
		b.buyTree.Descend(func(item limitTreeNode) bool {
			if item.LimitRef == nil || count == depth {
				return false
			}
			sn.Buys = append(sn.Buys, NewBookSnapshotRecord(item.LimitRef.Price, item.LimitRef.Size))
			count++
			return true
		})
	}(&wg, b, depth)
	wg.Add(1)
	go func(wg *sync.WaitGroup, b *book, depth int) {
		defer wg.Done()
		b.sellMu.RLock()
		defer b.sellMu.RUnlock()
		count := 0
		b.sellTree.Ascend(func(item limitTreeNode) bool {
			if item.LimitRef == nil || count == depth {
				return false
			}
			sn.Sells = append(sn.Sells, NewBookSnapshotRecord(item.LimitRef.Price, item.LimitRef.Size))
			count++
			return true
		})
	}(&wg, b, depth)
	wg.Wait()
	return sn
}
