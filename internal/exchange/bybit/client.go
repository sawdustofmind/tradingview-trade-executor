package bybit

import (
	"github.com/frenswifbenefits/myfren/internal/entity"
	bybit "github.com/wuhewuhe/bybit.go.api"
	"go.uber.org/zap"
)

type Client struct {
	logger         *zap.Logger
	customer       entity.Customer
	exchangeClient *bybit.Client
}
