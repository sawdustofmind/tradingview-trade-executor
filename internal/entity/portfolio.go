package entity

import (
	"encoding/json"

	"github.com/uptrace/bun"
)

const (
	StrategyTypeBase = "base"
	StrategyTypeDCA  = "dca"
)

type Portfolio struct {
	bun.BaseModel `bun:"table:portfolio"`

	Id                     int64           `bun:"id,pk,autoincrement"`
	Name                   string          `bun:"name"`
	Description            string          `bun:"description"`
	ImageBase64            string          `bun:"image_base64"`
	YearPnl                string          `bun:"year_pnl"`
	AvgDelay               string          `bun:"avg_delay"`
	RiskLevel              string          `bun:"risk_level"`
	StrategyType           string          `bun:"strategy_type"`
	DCALevels              int64           `bun:"dca_levels"`
	Leverage               int64           `bun:"leverage"`
	CycleInvestmentPercent string          `bun:"cycle_investment_percent"`
	Holdings               json.RawMessage `bun:"holdings"`
}

type Holding struct {
	Coin    string
	Percent string
}

func (p *Portfolio) GetHoldings() ([]Holding, error) {
	var holdings []Holding
	err := json.Unmarshal(p.Holdings, &holdings)
	return holdings, err
}

func (p *Portfolio) SetHoldings(holdings []Holding) error {
	var err error
	p.Holdings, err = json.Marshal(holdings)
	return err
}
