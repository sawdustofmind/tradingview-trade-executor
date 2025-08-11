package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frenswifbenefits/myfren/internal/app"
	"github.com/frenswifbenefits/myfren/internal/config"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	logger := &zap.Logger{}
	if cfg.Debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync() //nolint

	app, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal("cannot create app", zap.Error(err))
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	err = app.Start(ctx)
	if err != nil {
		logger.Fatal("cannot start app", zap.Error(err))
	}

	exitSignal := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)

	<-exitSignal
	logger.Info("Shutting down")
	cancelFn()

	err = app.Shutdown(5 * time.Second)
	if err != nil {
		logger.Error("failed to shut down", zap.Error(err))
	} else {
		logger.Info("Shut down was graceful")
	}
}
