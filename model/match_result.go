package model

type MatchResult struct {
	Trades        []Trade             `json:"trades"`
	Cancellations []OrderCancellation `json:"cancellations"`
}
