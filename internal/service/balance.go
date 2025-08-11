package service

import (
	"context"
	"github.com/frenswifbenefits/myfren/internal/daemons"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/exchange/bybit"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"go.uber.org/zap"
)

type BalanceService struct {
	cb     *bybit.ClientBuilder
	repo   *repository.Repository
	cp     *daemons.CustomerPool
	logger *zap.Logger
}

func NewBalanceService(logger *zap.Logger, repo *repository.Repository, cb *bybit.ClientBuilder, cp *daemons.CustomerPool) *BalanceService {
	return &BalanceService{
		repo:   repo,
		cb:     cb,
		cp:     cp,
		logger: logger.With(zap.String("module", "BalanceService")),
	}
}

func (bs *BalanceService) Balance(ctx context.Context, customer entity.Customer, isTest bool) (map[string]string, error) {
	exchangeClient, err := bs.cb.Build(customer, isTest)
	if err != nil {
		return nil, err
	}

	balances, err := exchangeClient.Balance(ctx)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func (bs *BalanceService) Healthcheck(ctx context.Context, customer entity.Customer, isTest bool) (bool, error) {
	exchangeClient, err := bs.cb.Build(customer, isTest)
	if err != nil {
		return false, err
	}

	_, err = exchangeClient.GetPosition(ctx, "BTCUSDT")
	if err != nil {
		return false, err
	}

	return true, nil
}
