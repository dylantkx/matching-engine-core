package orderbook

import "github.com/dylantkx/matching-engine-core/model"

type bookOrder struct {
	Order         *model.Order
	prevBookOrder *bookOrder
	nextBookOrder *bookOrder
}
