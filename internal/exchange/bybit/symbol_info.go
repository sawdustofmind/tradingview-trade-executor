package bybit

import (
	"context"
	"errors"
	"fmt"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/oliveagle/jsonpath"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s *Client) GetSymbolLotSize(
	ctx context.Context,
	symbol string,
) (dto.SymbolInfo, error) {
	params := map[string]interface{}{
		"category": "linear",
		"symbol":   symbol,
	}

	s.logger.Info("get si",
		zap.Any("params", params),
	)

	res, err := s.exchangeClient.NewUtaBybitServiceWithParams(params).GetInstrumentInfo(ctx)
	if err != nil {
		return dto.SymbolInfo{}, err
	}

	if res.RetCode != 0 {
		return dto.SymbolInfo{}, errors.New(res.RetMsg)
	}

	siRes, err := jsonpath.JsonPathLookup(res.Result, fmt.Sprintf("$.list[?(@.symbol == %s)][0].lotSizeFilter", symbol))
	if err != nil {
		return dto.SymbolInfo{}, err
	}
	si, ok := siRes.(map[string]interface{})
	if !ok {
		return dto.SymbolInfo{}, errors.New("si result is not map")
	}

	lotSize := fmt.Sprint(si["qtyStep"])
	minOrderQty := fmt.Sprint(si["minOrderQty"])

	s.logger.Info("get symbol info",
		zap.String("lotSize", lotSize),
		zap.String("minOrderQty", minOrderQty),
		zap.String("symbol", symbol),
	)

	minOrderQtyD, err := decimal.NewFromString(minOrderQty)
	if err != nil {
		return dto.SymbolInfo{}, err
	}
	tickSize, err := decimal.NewFromString(lotSize)
	if err != nil {
		return dto.SymbolInfo{}, err
	}
	return dto.SymbolInfo{
		MinOrderQty: minOrderQtyD,
		TickSize:    tickSize,
	}, nil
}
