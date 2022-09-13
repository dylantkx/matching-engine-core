package model

type OrderType = string

const (
	OrderType_Market OrderType = "MARKET"
	OrderType_Limit  OrderType = "LIMIT"
)
