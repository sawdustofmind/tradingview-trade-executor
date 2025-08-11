package dto

import "github.com/shopspring/decimal"

type OrderBookTop struct {
	Ask decimal.Decimal
	Bid decimal.Decimal
}
