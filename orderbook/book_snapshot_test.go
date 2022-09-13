package orderbook_test

import (
	"encoding/json"
	"testing"

	"github.com/dylantkx/matching-engine-core/orderbook"
	"github.com/shopspring/decimal"
)

func TestBookSnapshotMarshalJSONEmpty(t *testing.T) {
	sn := orderbook.NewBookSnapshot()

	b, err := json.Marshal(&sn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	j := `{"buys":[],"sells":[]}`
	if string(b) != j {
		t.Fatalf("wrong json output, want %s but got %s\n", j, b)
	}
}

func TestBookSnapshotUnmarshalJSONEmpty(t *testing.T) {
	sn := orderbook.NewBookSnapshot()

	input := `{"buys":[],"sells":[]}`
	err := json.Unmarshal([]byte(input), &sn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(sn.Buys) != 0 {
		t.Fatalf("expect buys length to be 0, but got %d", len(sn.Buys))
	}
	if len(sn.Sells) != 0 {
		t.Fatalf("expect sells length to be 0, but got %d", len(sn.Sells))
	}
}

func TestBookSnapshotMarshalJSON(t *testing.T) {
	sn := orderbook.NewBookSnapshot()
	sn.Buys = append(sn.Buys, orderbook.NewBookSnapshotRecord(decimal.NewFromFloat(100), decimal.NewFromFloat(1)))
	sn.Sells = append(sn.Sells, orderbook.NewBookSnapshotRecord(decimal.NewFromFloat(150), decimal.NewFromFloat(1)))

	b, err := json.Marshal(&sn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	j := `{"buys":[{"price":"100","size":"1"}],"sells":[{"price":"150","size":"1"}]}`
	if string(b) != j {
		t.Fatalf("wrong json output, want %s but got %s\n", j, b)
	}
}

func TestBookSnapshotUnmarshalJSON(t *testing.T) {
	sn := orderbook.NewBookSnapshot()

	input := `{"buys":[{"price":"100","size":"1"}],"sells":[{"price":"150","size":"1"}]}`
	err := json.Unmarshal([]byte(input), &sn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(sn.Buys) != 1 {
		t.Fatalf("expect buys length to be 1, but got %d", len(sn.Buys))
	}
	if len(sn.Sells) != 1 {
		t.Fatalf("expect sells length to be 1, but got %d", len(sn.Sells))
	}
	if !sn.Buys[0].Price.Equals(decimal.NewFromFloat(100)) {
		t.Fatalf("expect first buys price to be 100, but got %s", sn.Buys[0].Price)
	}
}
