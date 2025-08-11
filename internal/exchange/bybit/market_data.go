package bybit

import (
	"context"
	"errors"
	"fmt"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s *Client) GetSymbolPrice(ctx context.Context, symbol string) (dto.OrderBookTop, error) {
	res, err := s.exchangeClient.NewClassicalBybitServiceWithParams(map[string]interface{}{
		"category": "linear",
		"symbol":   symbol,
		"limit":    1,
	}).GetOrderBookInfo(ctx)
	if err != nil {
		return dto.OrderBookTop{}, err
	}
	if res.RetCode != 0 {
		return dto.OrderBookTop{}, errors.New(res.RetMsg)
	}

	respMInf, ok := res.Result.(map[string]interface{})
	if !ok {
		return dto.OrderBookTop{}, errors.New("respMInf is not map")
	}

	asks, asksOk := respMInf["a"].([]interface{})
	bids, bidsOk := respMInf["b"].([]interface{})

	ask := ""
	bid := ""

	if asksOk && len(asks) > 0 {
		ask0, ok := asks[0].([]interface{})
		if ok && len(ask0) > 0 {
			ask = fmt.Sprint(ask0[0])
		}
	}
	if bidsOk && len(bids) > 0 {
		bid0, ok := bids[0].([]interface{})
		if ok && len(bid0) > 0 {
			bid = fmt.Sprint(bid0[0])
		}
	}

	s.logger.Info("get market data",
		zap.Any("symbol", symbol),
		zap.Any("ask", ask),
		zap.Any("bid", bid),
	)

	asdD, err := decimal.NewFromString(ask)
	if err != nil {
		return dto.OrderBookTop{}, err
	}
	bidD, err := decimal.NewFromString(bid)
	if err != nil {
		return dto.OrderBookTop{}, err
	}

	return dto.OrderBookTop{
		Ask: asdD,
		Bid: bidD,
	}, nil
}
