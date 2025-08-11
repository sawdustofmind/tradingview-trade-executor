package bybit

import (
	"context"
	"errors"

	"github.com/oliveagle/jsonpath"
	"go.uber.org/zap"
)

func (s *Client) Balance(
	ctx context.Context,
) (map[string]string, error) {

	params := map[string]interface{}{
		"category": "linear",
	}

	s.logger.Info("getting balances list",
		zap.Any("params", params),
	)

	res, err := s.exchangeClient.NewUtaBybitServiceWithParams(params).GetSingleCoinsBalance(ctx)
	if err != nil {
		return nil, err
	}

	if res.RetCode != 0 {
		return nil, errors.New(res.RetMsg)
	}

	balancesRes, err := jsonpath.JsonPathLookup(res.Result, "$.balance")
	if err != nil {
		return nil, err
	}
	balances, ok := balancesRes.([]interface{})
	if !ok {
		return nil, errors.New("balances result is not map")
	}

	s.logger.Info("get raw balance list",
		zap.Any("raw", balances),
	)

	balancesResp := make(map[string]string)
	for _, balance := range balances {
		coin := balance.(map[string]interface{})["coin"].(string)
		walletBalance := balance.(map[string]interface{})["walletBalance"].(string)
		balancesResp[coin] = walletBalance
	}

	s.logger.Info("get balances list",
		zap.Any("balances", balances),
	)

	return balancesResp, nil
}
