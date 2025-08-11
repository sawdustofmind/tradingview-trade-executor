package dto

import "github.com/shopspring/decimal"

type SymbolInfo struct {
	MinOrderQty decimal.Decimal
	TickSize    decimal.Decimal
}
