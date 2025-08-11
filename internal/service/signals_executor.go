package service

import (
	"context"
	"fmt"
	"github.com/frenswifbenefits/myfren/internal/daemons"
	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/exchange/bybit"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	SideBuy  = "Buy"
	SideSell = "Sell"
)

type SignalsExecutor struct {
	cb     *bybit.ClientBuilder
	repo   *repository.Repository
	cp     *daemons.CustomerPool
	logger *zap.Logger
}

func NewSignalsExecutor(logger *zap.Logger, repo *repository.Repository, cb *bybit.ClientBuilder, cp *daemons.CustomerPool) *SignalsExecutor {
	return &SignalsExecutor{
		cb:     cb,
		cp:     cp,
		repo:   repo,
		logger: logger.With(zap.String("module", "signal_executor")),
	}
}

func (se *SignalsExecutor) ExecuteSignal(ctx context.Context, signal dto.Signal) {
	portfolio, err := se.repo.FindPortfolio(signal.StrategyName)
	if err != nil {
		se.logger.Error("cannot find portfolio", zap.Any("portfolio", signal), zap.Error(err))
		return
	}

	subs, err := se.repo.GetActivePortfolioSubscriptions(portfolio.Id, signal.Exchange)
	if err != nil {
		se.logger.Error("cannot get active subscriptions", zap.Any("portfolio", portfolio), zap.Error(err))
		return
	}

	se.logger.Info("found active subscriptions",
		zap.Any("subscriptions", subs),
		zap.Any("portfolio", portfolio),
		zap.Any("signal", signal))

	for _, sub := range subs {
		if sub.Exchange != "BYBIT" {
			se.logger.Warn("ignoring subscription", zap.Any("subscription", sub))
			continue
		}
		err = se.ExecuteStrategy(ctx, signal, portfolio, sub)
		if err != nil {
			se.logger.Error("cannot execute portfolio",
				zap.Any("portfolio", signal),
				zap.Any("sub", sub),
				zap.Any("portfolio", portfolio),
				zap.Error(err))
		}
	}
}

func (se *SignalsExecutor) ExecuteStrategy(
	ctx context.Context,
	signal dto.Signal,
	portfolio entity.Portfolio,
	sub entity.PortfolioSubscription,
) error {

	customer, ok := se.cp.GetByID(sub.CustomerId)
	if !ok {
		se.logger.Warn("cannot find customer with existing sub", zap.Any("customer", sub.CustomerId))
		return nil
	}

	exchangeClient, err := se.cb.Build(*customer, sub.IsTest)
	if err != nil {
		return err
	}

	shareMultiplier := decimal.NewFromFloat(1)
	coinFound := false
	holdings, err := portfolio.GetHoldings()
	if err != nil {
		return err
	}
	for _, holding := range holdings {
		if holding.Coin+"USDT" == signal.Symbol {
			shareMultiplier, err = decimal.NewFromString(holding.Percent)
			if err != nil {
				return err
			}
			shareMultiplier = shareMultiplier.Div(decimal.NewFromFloat(100))
			coinFound = true
			break
		}
	}
	if !coinFound {
		return fmt.Errorf("cannot find holdings for signal %s", signal.Symbol)
	}

	var details map[string]interface{}
	var corrId = uuid.New()
	var execErr error

	switch portfolio.StrategyType {
	case entity.StrategyTypeBase:
		details, execErr = se.execSimple(ctx, signal, portfolio, sub, exchangeClient, corrId, shareMultiplier)
	case entity.StrategyTypeDCA:
		details, execErr = se.execDca(ctx, signal, portfolio, sub, exchangeClient, corrId, shareMultiplier)
	default:
		return fmt.Errorf("unknown portfolio type: %s", portfolio.StrategyType)
	}

	if len(details) == 0 {
		details = map[string]interface{}{}
	}

	execErrStr := ""
	if execErr != nil {
		se.logger.Info("cannot execute portfolio", zap.Error(execErr))
		execErrStr = execErr.Error()
	}

	_, insertLogErr := se.repo.InsertAction(entity.Action{
		CorrId:            corrId,
		CustomerId:        customer.Id,
		SubId:             sub.Id,
		PortfolioId:       portfolio.Id,
		ActionType:        signal.Action,
		Details:           details,
		NeedToFetchTrades: execErr == nil,
		Error:             execErrStr,
	})

	if insertLogErr != nil {
		return insertLogErr
	}

	return nil
}

func (se *SignalsExecutor) execSimple(
	ctx context.Context,
	signal dto.Signal,
	portfolio entity.Portfolio,
	sub entity.PortfolioSubscription,
	exchangerClient *bybit.Client,
	corrId uuid.UUID,
	shareMultiplier decimal.Decimal,
) (map[string]interface{}, error) {
	switch signal.Action {
	case "open", "long", "buy", "entry", "short", "sell":
		side := SideBuy
		if signal.Action == "short" || signal.Action == "sell" {
			side = SideSell
		}

		positionSize, err := exchangerClient.GetPosition(ctx, signal.Symbol)
		if err != nil {
			return nil, err
		}

		if !positionSize.IsZero() {
			se.logger.Warn("position is not empty",
				zap.Any("position", positionSize),
				zap.Any("portfolio", portfolio),
				zap.Any("sub", sub),
			)
			return nil, nil
		}

		err = exchangerClient.SetLeverage(ctx, signal.Symbol, int(portfolio.Leverage))
		if err != nil {
			return nil, err
		}

		amountPercent, err := decimal.NewFromString(portfolio.CycleInvestmentPercent)
		if err != nil {
			return nil, err
		}

		amountPercent = amountPercent.Mul(shareMultiplier)

		quantity, err := getQuantityFromAmount(ctx, exchangerClient, sub.Amount, amountPercent, signal.Symbol, side)
		if err != nil {
			return nil, err
		}

		details := map[string]interface{}{
			"side":     side,
			"quantity": quantity.String(),
			"is_test":  sub.IsTest,
		}
		err = exchangerClient.PlaceOrder(ctx, signal.Symbol, side, quantity.String(), corrId.String())
		if err != nil {
			return details, err
		}

		return details, nil
	case "close":
		positionSize, err := exchangerClient.GetPosition(ctx, signal.Symbol)
		if err != nil {
			return nil, err
		}

		if positionSize.IsZero() {
			se.logger.Info("position is empty",
				zap.Any("position", positionSize),
				zap.Any("portfolio", portfolio),
				zap.Any("sub", sub),
			)
			return nil, nil
		}

		side := SideSell
		if positionSize.IsNegative() {
			side = SideBuy
			positionSize = positionSize.Neg()
		}

		details := map[string]interface{}{
			"side":     side,
			"quantity": positionSize.String(),
			"is_test":  sub.IsTest,
		}
		err = exchangerClient.PlaceOrder(ctx, signal.Symbol, side, positionSize.String(), corrId.String())
		if err != nil {
			return nil, err
		}

		return details, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", signal.Action)
	}
}

func (se *SignalsExecutor) execDca(
	ctx context.Context,
	signal dto.Signal,
	portfolio entity.Portfolio,
	sub entity.PortfolioSubscription,
	exchangerClient *bybit.Client,
	corrId uuid.UUID,
	shareMultiplier decimal.Decimal,
) (map[string]interface{}, error) {
	switch signal.Action {
	case "open", "long", "buy", "entry":
		side := SideBuy

		err := exchangerClient.SetLeverage(ctx, signal.Symbol, int(portfolio.Leverage))
		if err != nil {
			return nil, err
		}

		amountPercent, err := decimal.NewFromString(portfolio.CycleInvestmentPercent)
		if err != nil {
			return nil, err
		}

		amountPercent = amountPercent.Div(decimal.NewFromInt(portfolio.DCALevels)).Mul(shareMultiplier)

		quantity, err := getQuantityFromAmount(ctx, exchangerClient, sub.Amount, amountPercent, signal.Symbol, side)
		if err != nil {
			return nil, err
		}

		details := map[string]interface{}{
			"side":     side,
			"quantity": quantity.String(),
			"is_test":  sub.IsTest,
		}
		err = exchangerClient.PlaceOrder(ctx, signal.Symbol, side, quantity.String(), corrId.String())
		if err != nil {
			return details, err
		}

		return details, nil
	case "close":
		positionSize, err := exchangerClient.GetPosition(ctx, signal.Symbol)
		if err != nil {
			return nil, err
		}

		if positionSize.IsZero() {
			se.logger.Info("position is empty",
				zap.Any("position", positionSize),
				zap.Any("portfolio", portfolio),
				zap.Any("sub", sub),
			)
			return nil, nil
		}

		side := SideSell
		if positionSize.IsNegative() {
			side = SideBuy
			positionSize = positionSize.Neg()
		}

		details := map[string]interface{}{
			"side":     side,
			"quantity": positionSize.String(),
			"is_test":  sub.IsTest,
		}
		err = exchangerClient.PlaceOrder(ctx, signal.Symbol, side, positionSize.String(), corrId.String())
		if err != nil {
			return nil, err
		}

		return details, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", signal.Action)
	}
}

var hundred = decimal.NewFromInt(100)

func getQuantityFromAmount(
	ctx context.Context,
	exchangeClient *bybit.Client,
	amount string,
	amountPercent decimal.Decimal,
	symbol string,
	side string,
) (decimal.Decimal, error) {
	amountD, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Zero, err
	}

	priceInfo, err := exchangeClient.GetSymbolPrice(ctx, symbol)
	if err != nil {
		return decimal.Decimal{}, err
	}

	price := priceInfo.Ask
	if side != SideBuy {
		price = priceInfo.Bid
	}

	si, err := exchangeClient.GetSymbolLotSize(ctx, symbol)
	if err != nil {
		return decimal.Decimal{}, err
	}

	quantity := amountD.
		Mul(amountPercent).
		Div(hundred).
		Div(price).
		Div(si.TickSize).
		Truncate(0).
		Mul(si.TickSize)
	if quantity.LessThan(si.MinOrderQty) {
		requiredQuoteAmount := si.MinOrderQty.Mul(price)
		return decimal.Decimal{}, fmt.Errorf("order violates minimum amount '%s'<'%s'", amount, requiredQuoteAmount.String())
	}

	return quantity, nil
}
