package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type config struct {
	serverAddress string `env:"SERVER_ADDRESS"`
	baseURL       string `env:"BASE_URL"`
}

var cfg config

func init() {
	cfg = config{
		serverAddress: "localhost:8080",
		baseURL:       "http://localhost:8080/",
	}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func ServerAddress() string {
	return cfg.serverAddress
}

func BaseURL() string {
	return cfg.baseURL
}
