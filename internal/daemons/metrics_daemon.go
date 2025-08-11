package daemons

import (
	"strconv"
	"time"

	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/metrics"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"go.uber.org/zap"
)

type MetricsDaemon struct {
	logger *zap.Logger

	repo *repository.Repository
}

func NewMetricsDaemon(logger *zap.Logger, repo *repository.Repository) *MetricsDaemon {
	logger = logger.With(zap.String("component", "metrics_daemon"))
	return &MetricsDaemon{
		logger: logger,
		repo:   repo,
	}
}

func (cp *MetricsDaemon) Iteration() error {
	portfolios, err := cp.repo.GetPortfolios()
	if err != nil {
		return err
	}

	activeSubs, err := cp.repo.GetAllActiveSubscription()
	if err != nil {
		return err
	}
	portfolioMap := make(map[int64]entity.Portfolio)
	for _, portfolio := range portfolios {
		portfolioMap[portfolio.Id] = portfolio
	}

	metrics.SubscriptionsCount.Reset()
	metrics.SubscriptionsAmount.Reset()
	for _, activeSub := range activeSubs {
		portfolio, ok := portfolioMap[activeSub.PortfolioId]
		if !ok {
			continue
		}

		amount := 0.0
		amount, err = strconv.ParseFloat(activeSub.Amount, 64)
		if err != nil {
			amount = 0.0
			err = nil //nolint
		}
		subType := "main"
		if activeSub.IsTest {
			subType = "test"
		}
		metrics.SubscriptionsCount.WithLabelValues(portfolio.Name, subType).Inc()
		metrics.SubscriptionsAmount.WithLabelValues(portfolio.Name, subType).Add(amount)
	}
	metrics.StrategiesCount.Set(float64(len(portfolios)))
	return nil
}

func (cp *MetricsDaemon) Run(period time.Duration) error {
	err := cp.Iteration()
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for range ticker.C {
			err := cp.Iteration()
			if err != nil {
				cp.logger.Error("failed to iterate metrics daemon", zap.Error(err))
			}
		}
	}()

	return nil
}
