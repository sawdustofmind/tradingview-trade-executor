package bybit

import (
	"context"
	"errors"
	"strings"

	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/google/uuid"
	"github.com/oliveagle/jsonpath"
	"go.uber.org/zap"
)

func (s *Client) GetTrades(
	ctx context.Context,
	corrId uuid.UUID,
	subId int64,
	strategyId int64,
) ([]entity.ActionTrade, error) {
	logger := s.logger.With(
		zap.String("corrId", corrId.String()),
	)

	params := map[string]interface{}{
		"category":    "linear",
		"orderLinkId": corrId.String(),
	}

	s.logger.Info("getting trade list",
		zap.Any("params", params),
	)

	res, err := s.exchangeClient.NewUtaBybitServiceWithParams(params).GetTradeHistory(ctx)
	if err != nil {
		return nil, err
	}

	if res.RetCode != 0 {
		return nil, errors.New(res.RetMsg)
	}

	tradesRes, err := jsonpath.JsonPathLookup(res.Result, "$.list")
	if err != nil {
		return nil, err
	}
	trades, ok := tradesRes.([]interface{})
	if !ok {
		return nil, errors.New("trades result is not map")
	}

	logger.Info("get raw trade list",
		zap.Any("raw", trades),
	)

	tradeEntities := make([]entity.ActionTrade, 0, len(trades))
	for _, trade := range trades {
		mInt, ok := trade.(map[string]interface{})
		if !ok {
			return nil, errors.New("trades result is not map")
		}

		orderLinkId := mInt["orderLinkId"].(string)
		if orderLinkId != corrId.String() {
			logger.Warn("orderLinkId not equal corrId", zap.String("corrId", corrId.String()),
				zap.String("orderLinkId", orderLinkId))
			continue
		}
		side := mInt["side"].(string)
		if strings.ToLower(side) == "buy" {
			side = "buy"
		} else {
			side = "sell"
		}

		tradeEntities = append(tradeEntities, entity.ActionTrade{
			CorrId:      corrId,
			CustomerId:  s.customer.Id,
			Exchange:    "BYBIT",
			Side:        side,
			SubId:       subId,
			PortfolioId: strategyId,
			Symbol:      mInt["symbol"].(string),
			Quantity:    mInt["execQty"].(string),
			Price:       mInt["execPrice"].(string),
			Commission:  mInt["execFee"].(string),
		})
	}

	logger.Info("get trade list",
		zap.Any("trades", trades),
	)
	return tradeEntities, nil
}
