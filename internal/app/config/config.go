package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

var cfg config

func init() {
	cfg = config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func ServerAddress() string {
	return cfg.ServerAddress
}

func BaseURL() string {
	return cfg.BaseURL
}
