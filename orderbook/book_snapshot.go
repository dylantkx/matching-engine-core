package orderbook

type BookSnapshot struct {
	Buys  []*bookSnapshotRecord `json:"buys"`
	Sells []*bookSnapshotRecord `json:"sells"`
}

func NewBookSnapshot() *BookSnapshot {
	return &BookSnapshot{
		Buys:  make([]*bookSnapshotRecord, 0),
		Sells: make([]*bookSnapshotRecord, 0),
	}
}

func (sn *BookSnapshot) AppendBuy(record *bookSnapshotRecord) {
	sn.Buys = append(sn.Buys, record)
}

func (sn *BookSnapshot) AppendSell(record *bookSnapshotRecord) {
	sn.Sells = append(sn.Sells, record)
}
