package daemons

import (
	"context"
	"time"

	"github.com/frenswifbenefits/myfren/internal/exchange/bybit"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"go.uber.org/zap"
)

type TradesDaemon struct {
	logger *zap.Logger
	cb     *bybit.ClientBuilder
	cp     *CustomerPool
	repo   *repository.Repository
}

func NewTradesDaemon(logger *zap.Logger, repo *repository.Repository, cp *CustomerPool, cb *bybit.ClientBuilder) *TradesDaemon {
	logger = logger.With(zap.String("component", "trades_daemon"))
	return &TradesDaemon{
		logger: logger,
		repo:   repo,
		cp:     cp,
		cb:     cb,
	}
}

func (td *TradesDaemon) Iteration() error {
	actions, err := td.repo.GetUnprocessedActions()
	if err != nil {
		return err
	}

	if len(actions) == 0 {
		return nil
	}

	td.logger.Info("processing actions", zap.Int("count", len(actions)))
	for _, action := range actions {
		customer, ok := td.cp.GetByID(action.CustomerId)
		if !ok {
			td.logger.Warn("skip customer to fetch trades", zap.Any("action", action))
			continue
		}

		sub, err := td.repo.GetPortfolioSubscriptionById(action.SubId)
		if err != nil {
			return err
		}

		client, err := td.cb.Build(*customer, sub.IsTest)
		if err != nil {
			return err
		}
		trades, err := client.GetTrades(context.Background(), action.CorrId, action.SubId, action.PortfolioId)
		if err != nil {
			return err
		}

		for _, trade := range trades {
			_, err = td.repo.InsertActionTrade(trade)
			if err != nil {
				return err
			}
		}

		err = td.repo.MarkActionAsProcessed(int64(action.Id))
		if err != nil {
			return err
		}
	}
	return nil
}

func (td *TradesDaemon) Run(period time.Duration) error {
	err := td.Iteration()
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for range ticker.C {
			err := td.Iteration()
			if err != nil {
				td.logger.Error("failed to iterate trades daemon", zap.Error(err))
			}
		}
	}()

	return nil
}
