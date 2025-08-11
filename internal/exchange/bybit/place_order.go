package bybit

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

func (s *Client) PlaceOrder(
	ctx context.Context,
	symbol string,
	side string,
	quantity string,
	corrId string,
) error {
	params := map[string]interface{}{
		"category":    "linear",
		"symbol":      symbol,
		"side":        side,
		"orderType":   "Market",
		"qty":         quantity,
		"price":       "0",
		"timeInForce": "GTC",
		"orderLinkId": corrId,
	}

	s.logger.Info("order create",
		zap.Any("params", params),
	)

	res, err := s.exchangeClient.NewUtaBybitServiceWithParams(params).PlaceOrder(ctx)
	if err != nil {
		return err
	}

	if res.RetCode != 0 {
		return errors.New(res.RetMsg)
	}
	return nil
}
