package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dylantkx/matching-engine-core/model"
	"github.com/shopspring/decimal"
)

func TestMarshalTrade(t *testing.T) {
	tr := model.Trade{
		BuyOrderID:   1,
		SellOrderID:  2,
		Units:        decimal.NewFromFloat(1.5),
		Price:        decimal.NewFromFloat(100),
		IsBuyerMaker: false,
		EventTime:    model.Timestamp{Time: time.Unix(1663079295, 0)},
	}

	b, err := json.Marshal(&tr)
	if err != nil {
		t.Fatalf(err.Error())
	}
	out := `{"buyOrderId":1,"sellOrderId":2,"units":"1.5","price":"100","isBuyerMaker":false,"eventTime":1663079295}`
	if string(b) != out {
		t.Fatalf("wrong output, got %s", b)
	}
}
