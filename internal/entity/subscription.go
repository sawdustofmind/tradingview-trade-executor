package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type PortfolioSubscription struct {
	bun.BaseModel `bun:"table:portfolio_subscription"`

	Id          int64     `bun:"id,pk,autoincrement"`
	PortfolioId int64     `bun:"portfolio_id"`
	CustomerId  int64     `bun:"customer_id"`
	IsTest      bool      `bun:"is_test"`
	Exchange    string    `bun:"exchange"`
	Amount      string    `bun:"amount"`
	Status      string    `bun:"status"`
	Pnl         string    `bun:"pnl"`
	CreatedAt   time.Time `bun:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at"`
}
