package entity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"time"
)

type ActionTrade struct {
	bun.BaseModel `bun:"table:action_trades"`

	Id          int64     `bun:"id,pk,autoincrement"`
	CorrId      uuid.UUID `bun:"corr_id"`
	CustomerId  int64     `bun:"customer_id"`
	SubId       int64     `bun:"sub_id"`
	PortfolioId int64     `bun:"portfolio_id"`
	Exchange    string    `bun:"exchange"`
	Side        string    `bun:"side"`
	Symbol      string    `bun:"symbol"`
	Quantity    string    `bun:"quantity"`
	Price       string    `bun:"price"`
	Commission  string    `bun:"commission"`
	CreatedAt   time.Time `bun:"created_at"`
}
