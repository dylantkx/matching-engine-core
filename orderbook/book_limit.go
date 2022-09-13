package orderbook

import (
	"sync"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/shopspring/decimal"
)

type bookLimit struct {
	Price          decimal.Decimal
	Size           decimal.Decimal
	Volume         decimal.Decimal
	firstBookOrder *bookOrder
	lastBookOrder  *bookOrder
	bookOrderMap   map[int64]*bookOrder
	mu             sync.RWMutex
}

func NewBookLimit() *bookLimit {
	return &bookLimit{
		Price:        decimal.NewFromFloat(0),
		Size:         decimal.NewFromFloat(0),
		Volume:       decimal.NewFromFloat(0),
		bookOrderMap: make(map[int64]*bookOrder),
	}
}

func (bl *bookLimit) InsertOrUpdateOrder(order *model.Order) (isUpdate bool) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	o := bl.bookOrderMap[order.ID]
	if o != nil {
		bl.updateSize(o.Order.Units.Neg())
		o.Order = order
		bl.updateSize(o.Order.Units)
		isUpdate = true
		return
	}
	nbo := &bookOrder{
		Order: order,
	}
	if bl.firstBookOrder == nil {
		bl.firstBookOrder = nbo
		bl.lastBookOrder = nbo
		bl.Price = order.Price
	} else {
		bl.lastBookOrder.nextBookOrder = nbo
		nbo.prevBookOrder = bl.lastBookOrder
		bl.lastBookOrder = nbo
	}
	bl.updateSize(order.Units)
	bl.bookOrderMap[order.ID] = nbo
	return
}

func (bl *bookLimit) RemoveOrder(order *model.Order) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	o := bl.bookOrderMap[order.ID]
	if o == nil {
		return
	}
	if o.prevBookOrder != nil {
		o.prevBookOrder.nextBookOrder = o.nextBookOrder
	}
	if o.nextBookOrder != nil {
		o.nextBookOrder.prevBookOrder = o.prevBookOrder
	}
	if o == bl.firstBookOrder {
		bl.firstBookOrder = o.nextBookOrder
	}
	if o == bl.lastBookOrder {
		bl.lastBookOrder = o.prevBookOrder
	}
	bl.updateSize(o.Order.Units.Neg())
	delete(bl.bookOrderMap, order.ID)
}

func (bl *bookLimit) IsEmpty() bool {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	return bl.firstBookOrder == nil
}

func (bl *bookLimit) CountOrders() int {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	return len(bl.bookOrderMap)
}

func (bl *bookLimit) updateSize(change decimal.Decimal) {
	bl.Size = bl.Size.Add(change)
	if bl.Size.IsNegative() {
		bl.Size = decimal.NewFromFloat(0)
	}
	bl.Volume = bl.Size.Mul(bl.Price)
}
