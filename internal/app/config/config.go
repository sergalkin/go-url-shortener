// Package config - uses.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`  // server address without protocol
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"` // base URL
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`             // path to file with stored URLs when using memory mode
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`                  // dsn used to establish connection with database
}

// OptionConfig - callback that can be provided to NewConfig to construct config with non default params.
type OptionConfig func(*config)

var cfg config

func init() {
	cfg = config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// NewConfig - Generate new config based on provided options.
func NewConfig(opts ...OptionConfig) *config {
	for _, opt := range opts {
		opt(&cfg)
	}

	return &cfg
}

// WithServerAddress - Generate config with ServerAddress.
func WithServerAddress(addr string) OptionConfig {
	return func(c *config) {
		c.ServerAddress = addr
	}
}

// WithBaseURL - Generate config with BaseURL.
func WithBaseURL(baseURL string) OptionConfig {
	return func(c *config) {
		c.BaseURL = baseURL
	}
}

// WithFileStoragePath - Generate config with FileStoragePath.
func WithFileStoragePath(path string) OptionConfig {
	return func(c *config) {
		c.FileStoragePath = path
	}
}

// WithDatabaseConnection - Generate config with DatabaseDSN.
func WithDatabaseConnection(addr string) OptionConfig {
	return func(c *config) {
		c.DatabaseDSN = addr
	}
}

// ServerAddress - Get ServerAddress from config.
func ServerAddress() string {
	return cfg.ServerAddress
}

// BaseURL - Get BaseURL from config.
func BaseURL() string {
	return cfg.BaseURL
}

// FileStoragePath - get FileStoragePath From config.
func FileStoragePath() string {
	return cfg.FileStoragePath
}

// DatabaseDSN - get DatabaseDSN from config.
func DatabaseDSN() string {
	return cfg.DatabaseDSN
}
