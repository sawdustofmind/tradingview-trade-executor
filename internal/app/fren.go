package app

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/frenswifbenefits/myfren/internal/config"
	"github.com/frenswifbenefits/myfren/internal/crypt"
	"github.com/frenswifbenefits/myfren/internal/daemons"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/exchange/bybit"
	httpAdmin "github.com/frenswifbenefits/myfren/internal/http/admin"
	"github.com/frenswifbenefits/myfren/internal/http/middleware"
	httpUser "github.com/frenswifbenefits/myfren/internal/http/user"
	httpWebhook "github.com/frenswifbenefits/myfren/internal/http/webhook"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"github.com/frenswifbenefits/myfren/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

type App struct {
	cfg           *config.Config
	logger        *zap.Logger
	customerPool  *daemons.CustomerPool
	metricsDaemon *daemons.MetricsDaemon
	tradesDaemon  *daemons.TradesDaemon

	userServer       *http.Server
	adminServer      *http.Server
	webhookServer    *http.Server
	monitoringServer *http.Server
}

func NewApp(cfg *config.Config, logger *zap.Logger) (*App, error) {
	crypter := crypt.NewCrypter([]byte(cfg.Security.AESSalt))

	conn := pgdriver.NewConnector(pgdriver.WithDSN(cfg.DB.DSN))
	sqldb := sql.OpenDB(conn)
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}
	db := bun.NewDB(sqldb, pgdialect.New())
	db.RegisterModel(entity.M2M...)

	repo := repository.NewRepository(db, crypter)
	cb := bybit.NewClientBuilder(logger, cfg.Bybit)
	cp := daemons.NewCustomerPool(logger, repo)
	md := daemons.NewMetricsDaemon(logger, repo)
	td := daemons.NewTradesDaemon(logger, repo, cp, cb)
	se := service.NewSignalsExecutor(logger, repo, cb, cp)
	bs := service.NewBalanceService(logger, repo, cb, cp)

	userRouter := gin.Default()
	userRouter.Use(gin.Recovery(), middleware.PrometheusMiddleware())
	userServerImpl := httpUser.NewServerImpl(logger, repo, cp, bs, se)
	httpUser.RegisterHandlers(userRouter, userServerImpl)

	adminRouter := gin.Default()
	adminRouter.Use(gin.Recovery(), middleware.PrometheusMiddleware())
	adminServerImpl := httpAdmin.NewServerImpl(logger, repo, cp, se, cfg.AdminServer.Users)
	httpAdmin.RegisterHandlers(adminRouter, adminServerImpl)

	webhookRouter := gin.Default()
	webhookRouter.Use(gin.Recovery(), middleware.PrometheusMiddleware())
	webhookServerImpl := httpWebhook.NewServerImpl(logger, cfg.WebhookServer.TvWhitelistEnabled, se)
	httpWebhook.RegisterHandlers(webhookRouter, webhookServerImpl)

	userServer := &http.Server{
		Addr:    cfg.UserServer.Address,
		Handler: userRouter,
	}
	adminServer := &http.Server{
		Addr:    cfg.AdminServer.Address,
		Handler: adminRouter,
	}
	webhookServer := &http.Server{
		Addr:    cfg.WebhookServer.Address,
		Handler: webhookRouter,
	}
	monitoringServer := &http.Server{
		Addr:    cfg.Monitoring.Address,
		Handler: promhttp.Handler(),
	}
	return &App{
		cfg:              cfg,
		logger:           logger,
		customerPool:     cp,
		metricsDaemon:    md,
		tradesDaemon:     td,
		userServer:       userServer,
		adminServer:      adminServer,
		webhookServer:    webhookServer,
		monitoringServer: monitoringServer,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	err := app.customerPool.RunInvalidate(5 * time.Second)
	if err != nil {
		return err
	}

	err = app.metricsDaemon.Run(10 * time.Second)
	if err != nil {
		return err
	}

	err = app.tradesDaemon.Run(10 * time.Second)
	if err != nil {
		return err
	}

	go func() {
		app.logger.Info("user http server start listening", zap.String("address", app.userServer.Addr))
		err := app.userServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal("user http server closed unexpectedly", zap.Error(err))
		}
	}()
	go func() {
		err := app.adminServer.ListenAndServe()
		app.logger.Info("admin http server start listening", zap.String("address", app.adminServer.Addr))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal("admin http server closed unexpectedly", zap.Error(err))
		}
	}()
	go func() {
		app.logger.Info("webhook http server start listening", zap.String("address", app.webhookServer.Addr))
		err := app.webhookServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal("webhook http server closed unexpectedly", zap.Error(err))
		}
	}()
	go func() {
		app.logger.Info("monitoring http server start listening", zap.String("address", app.monitoringServer.Addr))
		err := app.monitoringServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal("monitoring http server closed unexpectedly", zap.Error(err))
		}
	}()
	return nil
}

func (app *App) Shutdown(deadline time.Duration) error {
	return nil
}
