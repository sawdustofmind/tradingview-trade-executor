package bybit

import (
	"errors"

	"github.com/frenswifbenefits/myfren/internal/config"
	"github.com/frenswifbenefits/myfren/internal/entity"
	bybit "github.com/wuhewuhe/bybit.go.api"
	"go.uber.org/zap"
)

type ClientBuilder struct {
	conf   *config.BybitConfig
	logger *zap.Logger
}

func NewClientBuilder(logger *zap.Logger, conf *config.BybitConfig) *ClientBuilder {
	return &ClientBuilder{
		conf:   conf,
		logger: logger,
	}
}

func (cb *ClientBuilder) Build(customer entity.Customer, isTest bool) (*Client, error) {
	apiKey, apiSecret, err := readBybitApiCredentials(customer, isTest)
	if err != nil {
		return nil, err
	}

	host := bybit.WithBaseURL(cb.getApiUrl(isTest))
	exchangeClient := bybit.NewBybitHttpClient(apiKey, apiSecret, host)

	return &Client{
		logger:         cb.logger.With(zap.Any("customer", customer), zap.Any("isTest", isTest)),
		customer:       customer,
		exchangeClient: exchangeClient,
	}, nil
}

func (cb *ClientBuilder) getApiUrl(isTest bool) string {
	if isTest {
		return cb.conf.TestRestApi
	}
	return cb.conf.MainRestApi
}

func readBybitApiCredentials(customer entity.Customer, isTest bool) (string, string, error) {
	if isTest {
		if customer.BybitTestApiKey == nil {
			return "", "", errors.New("bybit_test_api_key is required")
		}

		if customer.BybitTestApiSecret == nil {
			return "", "", errors.New("bybit_test_api_secret is required")
		}
		return *customer.BybitTestApiKey, *customer.BybitTestApiSecret, nil
	}

	if customer.BybitApiKey == nil {
		return "", "", errors.New("bybit_api_key is required")
	}

	if customer.BybitApiSecret == nil {
		return "", "", errors.New("bybit_api_secret is required")
	}
	return *customer.BybitApiKey, *customer.BybitApiSecret, nil
}
