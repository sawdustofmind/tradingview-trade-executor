package adminhttp

import (
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/http/converters"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) DeleteV1AdminPortfolio(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	req := &api_types.DeleteV1AdminPortfolioJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	err = s.repository.DeletePortfolio(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())

}

func (s *ServerImpl) PostV1AdminPortfolio(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	req := &api_types.PostV1AdminPortfolioJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	portfolio := entity.Portfolio{
		AvgDelay:               req.AvgDelay,
		CycleInvestmentPercent: req.CycleInvestmentPercent,
		DCALevels:              int64(req.DcaLevels),
		Description:            req.Description,
		ImageBase64:            req.ImageBase64,
		Leverage:               int64(req.Leverage),
		Name:                   req.Name,
		RiskLevel:              req.RiskLevel,
		StrategyType:           req.StrategyType,
		YearPnl:                req.YearPnl,
	}
	holdings := make([]entity.Holding, 0, len(req.Holdings))
	for _, h := range req.Holdings {
		holdings = append(holdings, entity.Holding{
			Coin:    h.Coin,
			Percent: h.Percent,
		})
	}
	err = portfolio.SetHoldings(holdings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	id, err := s.repository.InsertPortfolio(portfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	resp := api_types.IdResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}
	resp.Data.Id = id
	c.JSON(http.StatusOK, resp)

}

func (s *ServerImpl) PutV1AdminPortfolio(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	req := &api_types.PutV1AdminPortfolioJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	portfolio, err := s.repository.GetPortfolioById(req.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeErrorAPIResponse(err))
		return
	}

	if req.Name != nil {
		portfolio.Name = *req.Name
	}
	if req.Description != nil {
		portfolio.Description = *req.Description
	}
	if req.ImageBase64 != nil {
		portfolio.ImageBase64 = *req.ImageBase64
	}
	if req.YearPnl != nil {
		portfolio.YearPnl = *req.YearPnl
	}
	if req.AvgDelay != nil {
		portfolio.AvgDelay = *req.AvgDelay
	}
	if req.RiskLevel != nil {
		portfolio.RiskLevel = *req.RiskLevel
	}
	if req.StrategyType != nil {
		portfolio.StrategyType = *req.StrategyType
	}
	if req.DcaLevels != nil {
		portfolio.DCALevels = int64(*req.DcaLevels)
	}
	if req.Leverage != nil {
		portfolio.Leverage = int64(*req.Leverage)
	}
	if req.CycleInvestmentPercent != nil {
		portfolio.CycleInvestmentPercent = *req.CycleInvestmentPercent
	}
	if req.Holdings != nil {
		holdings := make([]entity.Holding, 0)
		for _, h := range *req.Holdings {
			holdings = append(holdings, entity.Holding{
				Coin:    h.Coin,
				Percent: h.Percent,
			})
		}
		err = portfolio.SetHoldings(holdings)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
			return
		}
	}

	err = s.repository.UpdatePortfolio(portfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())

}

func (s *ServerImpl) GetV1AdminPortfolioList(c *gin.Context) {
	if !s.auth(c) {
		return
	}

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
