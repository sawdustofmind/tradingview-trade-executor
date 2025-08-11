package bybit

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s *Client) GetPosition(
	ctx context.Context,
	symbol string,
) (decimal.Decimal, error) {
	params := map[string]interface{}{
		"category": "linear",
		"symbol":   symbol,
	}

	s.logger.Info("get positions",
		zap.Any("params", params),
	)

	res, err := s.exchangeClient.NewUtaBybitServiceWithParams(params).GetPositionList(ctx)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if res.RetCode != 0 {
		return decimal.Decimal{}, errors.New(res.RetMsg)
	}

	positionRes, err := jsonpath.JsonPathLookup(res.Result, "$.list[0]")
	if err != nil {
		return decimal.Decimal{}, err
	}
	position, ok := positionRes.(map[string]interface{})
	if !ok {
		return decimal.Decimal{}, errors.New("position result is not map")
	}

	size := fmt.Sprint(position["size"])
	side := fmt.Sprint(position["side"])

	sizeD, err := decimal.NewFromString(size)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if strings.ToLower(side) == "sell" {
		return sizeD.Neg(), nil
	}
	return sizeD, nil
}
