package userhttp

import (
	"github.com/frenswifbenefits/myfren/internal/daemons"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"github.com/frenswifbenefits/myfren/internal/service"
	"go.uber.org/zap"
)

var _ ServerInterface = (*ServerImpl)(nil)

type ServerImpl struct {
	repository *repository.Repository
	cp         *daemons.CustomerPool
	bs         *service.BalanceService
	se         *service.SignalsExecutor

	logger *zap.Logger
}

func NewServerImpl(
	logger *zap.Logger,
	repository *repository.Repository,
	cp *daemons.CustomerPool,
	bs *service.BalanceService,
	se *service.SignalsExecutor,

) *ServerImpl {
	return &ServerImpl{
		repository: repository,
		cp:         cp,
		bs:         bs,
		se:         se,
		logger:     logger.With(zap.String("service", "user_server")),
	}
}
