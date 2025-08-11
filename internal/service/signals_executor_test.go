package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/frenswifbenefits/myfren/internal/config"
	"github.com/frenswifbenefits/myfren/internal/crypt"
	"github.com/frenswifbenefits/myfren/internal/daemons"
	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/exchange/bybit"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

func TestDCA(t *testing.T) {
	// setup app
	cfg, err := config.ReadConfigWithPath("../../configs")
	require.NoError(t, err)

	logger, _ := zap.NewDevelopment()
	cb := bybit.NewClientBuilder(logger, cfg.Bybit)
	crypter := crypt.NewCrypter([]byte(cfg.Security.AESSalt))

	pgconn := pgdriver.NewConnector(pgdriver.WithDSN(cfg.DB.DSN))
	sqldb := sql.OpenDB(pgconn)
	err = sqldb.Ping()
	require.NoError(t, err)

	db := bun.NewDB(sqldb, pgdialect.New())
	db.RegisterModel(entity.M2M...)
	repo := repository.NewRepository(db, crypter)
	cp := daemons.NewCustomerPool(logger, repo)
	se := NewSignalsExecutor(logger, repo, cb, cp)
	bs := NewBalanceService(logger, repo, cb, cp)

	fixture(t, repo, cp)
	err = cp.Invalidate()
	require.NoError(t, err)

	customer := cp.GetAll()[0]
	balances, err := bs.Healthcheck(context.Background(), customer, true)
	require.NoError(t, err)
	fmt.Println(balances)

	//portfolios, err := repo.GetPortfolios()
	//require.NoError(t, err)

	se.ExecuteSignal(context.Background(), dto.Signal{
		Exchange:     "BYBIT",
		Symbol:       "XRPUSDT",
		StrategyName: "DCA",
		Action:       "buy",
	})

	time.Sleep(1 * time.Second)

	se.ExecuteSignal(context.Background(), dto.Signal{
		Exchange:     "BYBIT",
		Symbol:       "XRPUSDT",
		StrategyName: "DCA",
		Action:       "close",
	})
}

func fixture(t *testing.T, repo *repository.Repository, cp *daemons.CustomerPool) {
	apiKey := "ul441O3bktOCwca4Vk"
	apiSecret := "3NinztO1CTIxcHKpT8W5tcVdoqDYEvfXLioH"

	customer := entity.Customer{
		Id:                 1,
		Username:           "testuser",
		BybitTestApiKey:    &apiKey,
		BybitTestApiSecret: &apiSecret,
	}
	_, err := repo.InsertCustomer(customer)
	if err != nil && strings.Contains(err.Error(), "violates unique constraint") {
		err = nil
	}
	require.NoError(t, err)

	err = repo.UpdateCustomer(customer)
	require.NoError(t, err)

	err = cp.Invalidate()
	require.NoError(t, err)

	holdings, err := json.Marshal([]entity.Holding{
		{
			Coin:    "XRP",
			Percent: "0.4",
		},
		{
			Coin:    "JUP",
			Percent: "0.6",
		},
	})
	require.NoError(t, err)

	_, err = repo.InsertPortfolio(entity.Portfolio{
		Name:                   "DCA",
		Description:            "DCA WIP",
		YearPnl:                "47",
		AvgDelay:               "Instant",
		RiskLevel:              "Medium",
		StrategyType:           "dca",
		DCALevels:              7,
		Leverage:               2,
		CycleInvestmentPercent: "10000",
		Holdings:               holdings,
	})
	if err != nil && strings.Contains(err.Error(), "violates unique constraint") {
		err = nil
	}
	require.NoError(t, err)

	portfolios, err := repo.GetPortfolios()
	require.NoError(t, err)

	subs, err := repo.GetPortfolioSubscription(cp.GetAll()[0].Id)
	require.NoError(t, err)
	if len(subs) != 0 {
		return
	}

	_, err = repo.InsertPortfolioSubscription(entity.PortfolioSubscription{
		Id:          1,
		PortfolioId: portfolios[0].Id,
		CustomerId:  cp.GetAll()[0].Id,
		IsTest:      true,
		Amount:      "100",
		Status:      "active",
		Pnl:         "0",
		Exchange:    "BYBIT",
	})
	require.NoError(t, err)
}
