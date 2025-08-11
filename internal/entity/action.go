package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Action struct {
	bun.BaseModel `bun:"table:action"`

	Id                int64                  `bun:"id,pk,autoincrement"`
	CorrId            uuid.UUID              `bun:"corr_id"`
	CustomerId        int64                  `bun:"customer_id"`
	SubId             int64                  `bun:"sub_id"`
	PortfolioId       int64                  `bun:"portfolio_id"`
	ActionType        string                 `bun:"action_type"`
	Details           map[string]interface{} `bun:"details"`
	NeedToFetchTrades bool                   `bun:"need_to_fetch_trades"`
	Error             string                 `bun:"error"`
	CreatedAt         time.Time              `bun:"created_at"`
}
