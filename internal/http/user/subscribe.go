package userhttp

import (
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/entity"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) GetV1UserPortfolioSubscriptionsList(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	strategies, err := s.repository.GetPortfolioSubscription(int64(customer.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeBadRequestAPIResponse(err))
		return
	}
	resp := api_types.SubscriptionListResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}

	for _, subscription := range strategies {
		resp.Data = append(resp.Data, api_types.Subscription{
			Amount:      subscription.Amount,
			CreatedAt:   subscription.CreatedAt.String(),
			Id:          subscription.Id,
			IsTest:      subscription.IsTest,
			Pnl:         subscription.Pnl,
			PortfolioId: subscription.PortfolioId,
			Status:      subscription.Status,
			UpdatedAt:   subscription.UpdatedAt.String(),
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (s *ServerImpl) PostV1UserPortfolioSubscribe(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	req := &api_types.PostV1UserPortfolioSubscribeJSONRequestBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	healthy, err := s.bs.Healthcheck(c, *customer, req.IsTest)
	if err != nil || !healthy {
		c.JSON(http.StatusBadRequest, dto.MakeBadApiKeyAPIResponse(err))
		return
	}

	id, err := s.repository.InsertPortfolioSubscription(entity.PortfolioSubscription{
		PortfolioId: req.PortfolioId,
		CustomerId:  customer.Id,
		IsTest:      req.IsTest,
		Amount:      req.Amount,
		Status:      "active",
		Pnl:         "0",
		Exchange:    "BYBIT",
	})
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

func (s *ServerImpl) PostV1UserPortfolioUnsubscribe(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	req := &api_types.PostV1UserPortfolioUnsubscribeJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	sub, err := s.repository.GetPortfolioSubscriptionById(req.SubscriptionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}
	portfolio, err := s.repository.GetPortfolioById(sub.PortfolioId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	symbols, err := portfolio.GetHoldings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}
	for _, symbol := range symbols {
		err = s.se.ExecuteStrategy(c, dto.Signal{
			Exchange:     sub.Exchange,
			Symbol:       symbol.Coin + "USDT",
			StrategyName: portfolio.Name,
			Action:       "close",
		}, portfolio, sub)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
			return
		}
	}

	err = s.repository.PortfolioUnsubscribe(customer.Id, req.SubscriptionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeBadRequestAPIResponse(err))
		return
	}
	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())
}

func (s *ServerImpl) GetV1UserActionsList(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	actions, err := s.repository.GetActions(int64(customer.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeBadRequestAPIResponse(err))
		return
	}
	resp := api_types.ActionListResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}

	for _, action := range actions {
		//details := make(map[string]interface{})
		//if len(action.Details) > 0 {
		//	_ = json.Unmarshal(action.Details, &details)
		//}
		resp.Data = append(resp.Data, api_types.Action{
			ActionType:  action.ActionType,
			CorrId:      action.CorrId.String(),
			CreatedAt:   action.CreatedAt.String(),
			Details:     action.Details,
			Error:       action.Error,
			Id:          action.Id,
			PortfolioId: action.PortfolioId,
			SubId:       action.SubId,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (s *ServerImpl) GetV1UserTradesList(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	trades, err := s.repository.GetActionTrades(int64(customer.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeBadRequestAPIResponse(err))
		return
	}
	resp := api_types.TradeListResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}

	for _, trade := range trades {
		resp.Data = append(resp.Data, api_types.Trade{
			Commission:  trade.Commission,
			CorrId:      trade.CorrId.String(),
			CreatedAt:   trade.CreatedAt.String(),
			Exchange:    trade.Exchange,
			Id:          trade.Id,
			Price:       trade.Price,
			Quantity:    trade.Quantity,
			Side:        trade.Side,
			PortfolioId: trade.PortfolioId,
			SubId:       trade.SubId,
			Symbol:      trade.Symbol,
		})
	}

	c.JSON(http.StatusOK, resp)
}
