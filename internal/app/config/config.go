// Package config - uses.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"` // server address without protocol
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`      // base URL
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"" json:"file_storage_path"`         // path to file with stored URLs when using memory mode
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`                   // dsn used to establish connection with database
	JSONConfigPath  string `env:"CONFIG" envDefault:""`                                             // a path to config file
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" envDefault:"" json:"enable_https"`                   // a value used to determine http or https server will be run
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

// WithEnableHTTPS - Generate config with EnableHTTPS.
func WithEnableHTTPS(isEnabled bool) OptionConfig {
	return func(c *config) {
		c.EnableHTTPS = isEnabled
	}
}

// WithJSONConfig - Generate config with path to json.config file
func WithJSONConfig(path string) OptionConfig {
	return func(c *config) {
		c.JSONConfigPath = path
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

// EnableHTTPS - get EnableHTTPS from config.
func EnableHTTPS() bool {
	return cfg.EnableHTTPS
}

// JSONConfigPath - get path to config.json file
func JSONConfigPath() string {
	return cfg.JSONConfigPath
}

// SetJSONValues - set config zero values to json.config values
func (c *config) SetJSONValues() {
	// Open jsonFile
	jsonFile, err := os.Open(c.JSONConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// iterating over config struct via reflection to check for zero values,
	// if it has at least one field with none zero value, we do nothing, otherwise we will unmarshal data from
	// config.json file to config struct
	v := reflect.ValueOf(*c)
	isAllFieldsIsZeroValue := true
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name != "JSONConfigPath" && !v.Field(i).IsZero() {
			isAllFieldsIsZeroValue = false
			break
		}
	}

	if isAllFieldsIsZeroValue {
		json.Unmarshal(byteValue, &c)
	}
}
