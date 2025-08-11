package converters

import (
	"github.com/frenswifbenefits/myfren/internal/entity"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
)

func ConvertFren(fren entity.Fren) (api_types.Fren, error) {
	apiFren := api_types.Fren{
		Description: fren.Description,
		Id:          fren.Id,
		ImageBase64: fren.ImageBase64,
		Name:        fren.Name,
	}

	for _, portfolio := range fren.Portfolios {
		apiPortfolio, err := ConvertPortfolio(portfolio)
		if err != nil {
			return api_types.Fren{}, err
		}
		apiFren.Portfolios = append(apiFren.Portfolios, apiPortfolio)
	}
	return apiFren, nil
}
