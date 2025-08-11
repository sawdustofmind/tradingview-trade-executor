package converters

import (
	"fmt"
	"strings"

	"github.com/frenswifbenefits/myfren/internal/entity"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
)

func ConvertPortfolio(portfolio entity.Portfolio) (api_types.Portfolio, error) {

	holdings, err := portfolio.GetHoldings()
	if err != nil {
		return api_types.Portfolio{}, err
	}

	apiPortfolio := api_types.Portfolio{
		AvgDelay:               portfolio.AvgDelay,
		CycleInvestmentPercent: portfolio.CycleInvestmentPercent,
		DcaLevels:              int(portfolio.DCALevels),
		Description:            portfolio.Description,
		Id:                     portfolio.Id,
		ImageBase64:            portfolio.ImageBase64,
		Leverage:               int(portfolio.Leverage),
		Name:                   portfolio.Name,
		RiskLevel:              portfolio.RiskLevel,
		StrategyType:           portfolio.StrategyType,
		YearPnl:                portfolio.YearPnl,
	}

	for _, holding := range holdings {
		coinUrl := fmt.Sprintf(`https://assets.kraken.com/marketing/web/icons-uni-webp/s_%s.webp?i=kds`, strings.ToLower(holding.Coin))

		apiPortfolio.Holdings = append(apiPortfolio.Holdings, api_types.Holding{
			Coin:    holding.Coin,
			CoinImg: coinUrl,
			Percent: holding.Percent,
		})
	}

	return apiPortfolio, nil
}
