package model

import (
	"strconv"
	"time"
)

type Timestamp struct {
	Time time.Time
}

func (t *Timestamp) Unix() int64 {
	return t.Time.Unix()
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	s := strconv.Itoa(int(t.Unix()))
	return []byte(s), nil
}
