package entity

import (
	"github.com/uptrace/bun"
)

type Fren struct {
	bun.BaseModel `bun:"table:fren"`

	Id          int64       `bun:"id,pk,autoincrement"`
	Name        string      `bun:"name"`
	Description string      `bun:"description"`
	ImageBase64 string      `bun:"image_base64"`
	Portfolios  []Portfolio `bun:"m2m:fren_portfolio,join:Fren=Portfolio"`
}

type FrenPortfolio struct {
	bun.BaseModel `bun:"table:fren_portfolio"`
	FrenID        int64      `bun:"fren_id,pk"`
	Fren          *Fren      `bun:"rel:belongs-to,join:fren_id=id"`
	PortfolioId   int64      `bun:"portfolio_id,pk"`
	Portfolio     *Portfolio `bun:"rel:belongs-to,join:portfolio_id=id"`
}
