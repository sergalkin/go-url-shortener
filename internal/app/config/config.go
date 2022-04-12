package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

type OptionConfig func(*config)

var cfg config

func init() {
	cfg = config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func NewConfig(opts ...OptionConfig) *config {
	for _, opt := range opts {
		opt(&cfg)
	}

	return &cfg
}

func WithServerAddress(addr string) OptionConfig {
	return func(c *config) {
		c.ServerAddress = addr
	}
}

func WithBaseURL(baseURL string) OptionConfig {
	return func(c *config) {
		c.BaseURL = baseURL
	}
}

func WithFileStoragePath(path string) OptionConfig {
	return func(c *config) {
		c.FileStoragePath = path
	}
}

func ServerAddress() string {
	return cfg.ServerAddress
}

func BaseURL() string {
	return cfg.BaseURL
}

func FileStoragePath() string {
	return cfg.FileStoragePath
}
