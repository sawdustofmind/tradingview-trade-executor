package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	UserServer    *ServerConfig        `mapstructure:"user_server"`
	AdminServer   *AdminServerConfig   `mapstructure:"admin_server"`
	WebhookServer *WebhookServerConfig `mapstructure:"webhook_server"`
	Monitoring    *ServerConfig        `mapstructure:"monitoring_server"`
	DB            *DBConfig            `mapstructure:"db"`
	Security      *SecurityConfig      `mapstructure:"security"`
	Bybit         *BybitConfig         `mapstructure:"bybit"`
	Debug         bool                 `mapstructure:"debug"`
}

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

type AdminServerConfig struct {
	Address string            `mapstructure:"address"`
	Users   []AdminUserConfig `mapstructure:"users"`
}

type AdminUserConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"` // bcrytp hash
}

type WebhookServerConfig struct {
	Address            string `mapstructure:"address"`
	TvWhitelistEnabled bool   `mapstructure:"tv_whitelist_enabled"`
}

type RateLimiterConfig struct {
	RateLimit         int           `mapstructure:"rate_limit"`
	RateLimitInterval time.Duration `mapstructure:"rate_limit_interval"`
}

type SecurityConfig struct {
	AESSalt string `mapstructure:"aes_salt"`
}

type BybitConfig struct {
	MainRestApi string `mapstructure:"main_rest_api"`
	TestRestApi string `mapstructure:"test_rest_api"`
}

func ReadConfig() (*Config, error) {
	return ReadConfigWithPath("./configs")
}

func ReadConfigWithPath(path string) (*Config, error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	config := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
