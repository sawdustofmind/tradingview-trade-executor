package userhttp

import (
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/http/converters"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) GetV1UserPortfolioList(c *gin.Context) {
	portfolios, err := s.repository.GetPortfolios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	resp := api_types.PortfolioListResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}

	for _, portfolio := range portfolios {
		apiPortfolio, err := converters.ConvertPortfolio(portfolio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
			return
		}
		resp.Data = append(resp.Data, apiPortfolio)
	}

	c.JSON(http.StatusOK, resp)
}
