package adminhttp

import (
	"github.com/frenswifbenefits/myfren/internal/service"
	"sync"

	"github.com/frenswifbenefits/myfren/internal/config"
	"github.com/frenswifbenefits/myfren/internal/daemons"
	"github.com/frenswifbenefits/myfren/internal/repository"
	"go.uber.org/zap"
)

type ServerImpl struct {
	repository   *repository.Repository
	cp           *daemons.CustomerPool
	se           *service.SignalsExecutor
	logger       *zap.Logger
	authConf     []config.AdminUserConfig
	authTokensMu sync.RWMutex
	authTokens   map[string]struct{}
}

var _ ServerInterface = (*ServerImpl)(nil)

func NewServerImpl(
	logger *zap.Logger,
	repository *repository.Repository,
	cp *daemons.CustomerPool,
	se *service.SignalsExecutor,
	authConf []config.AdminUserConfig,
) *ServerImpl {
	return &ServerImpl{
		repository: repository,
		cp:         cp,
		authTokens: make(map[string]struct{}),
		authConf:   authConf,
		se:         se,
		logger:     logger.With(zap.String("service", "admin_server")),
	}
}
