package webhookhttp

import (
	"github.com/frenswifbenefits/myfren/internal/service"
	"go.uber.org/zap"
)

var _ ServerInterface = (*ServerImpl)(nil)

type ServerImpl struct {
	logger             *zap.Logger
	tvWhitelistEnabled bool

	si *service.SignalsExecutor
}

func NewServerImpl(
	logger *zap.Logger,
	tvWhitelistEnabled bool,
	si *service.SignalsExecutor,
) *ServerImpl {
	return &ServerImpl{
		tvWhitelistEnabled: tvWhitelistEnabled,
		si:                 si,
		logger:             logger.With(zap.String("service", "webhook_server")),
	}
}
