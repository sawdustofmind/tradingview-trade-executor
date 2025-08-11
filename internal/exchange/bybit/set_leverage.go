package bybit

import (
	"context"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

func (s *Client) SetLeverage(
	ctx context.Context,
	symbol string,
	leverage int,
) error {

	params := map[string]interface{}{
		"category":     "linear",
		"symbol":       symbol,
		"buyLeverage":  strconv.Itoa(leverage),
		"sellLeverage": strconv.Itoa(leverage),
	}

	s.logger.Info("set leverage",
		zap.String("symbol", symbol),
		zap.Int("leverage", leverage),
		zap.Any("params", params),
	)

	res, err := s.exchangeClient.NewUtaBybitServiceWithParams(params).SetPositionLeverage(ctx)
	if err != nil {
		return err
	}

	if res.RetCode != 0 &&
		res.RetCode != 110043 { // leverage not modified is ok error here
		return errors.New(res.RetMsg)
	}
	return nil
}
