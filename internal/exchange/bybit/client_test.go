package bybit

import (
	"context"
	"fmt"
	"testing"

	"github.com/frenswifbenefits/myfren/internal/config"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	cfg = &config.BybitConfig{
		TestRestApi: "https://api-demo.bybit.com",
	}

	apiKey    = "ul441O3bktOCwca4Vk"
	apiSecret = "3NinztO1CTIxcHKpT8W5tcVdoqDYEvfXLioH"

	customer = entity.Customer{
		Id:                 512,
		Username:           "testuser",
		BybitTestApiKey:    &apiKey,
		BybitTestApiSecret: &apiSecret,
	}
)

func TestFetchTrades(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cb := NewClientBuilder(logger, cfg)
	client, err := cb.Build(customer, true)
	require.NoError(t, err)

	trades, err := client.GetTrades(context.Background(), uuid.MustParse("d8b049b2-228b-451b-955b-dd5720eaf825"), 1, 2)
	require.NoError(t, err)

	logger.Info("got trades", zap.Any("trades", trades))
}

func TestGetSI(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cb := NewClientBuilder(logger, cfg)
	client, err := cb.Build(customer, true)
	require.NoError(t, err)

	for _, symbol := range []string{
		"BTCUSDT",
		"SOLUSDT",
		"AVAXUSDT",
		"NEARUSDT",
		"SUIUSDT",
		"XRPUSDT",
		"HBARUSDT",
		"RENDERUSDT",
		"JUPUSDT",
		"BALUSDT",
		"SHIBUSDT",
		"1INCHUSDT",
		"JUPUSDT",
	} {
		si, err := client.GetSymbolLotSize(context.Background(), symbol)
		if err != nil {
			fmt.Println("$$$", symbol, "not listed")
			continue
		}

		obTop, err := client.GetSymbolPrice(context.Background(), symbol)
		require.NoError(t, err)
		minQty := obTop.Ask.Mul(si.MinOrderQty)
		require.NoError(t, err)

		fmt.Println("$$$", symbol,
			"minUSDT", minQty,
			"minBTC", si.MinOrderQty.String(),
			"lotBTC", si.TickSize.String(),
			"price", obTop.Ask.String())
	}
}
