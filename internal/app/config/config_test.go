package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseURL(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "BaseURL can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewConfig()
			assert.NotEmpty(t, BaseURL())
		})
	}
}

func TestDatabaseDSN(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "DatabaseDSN can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithDatabaseConnection("test"))
			assert.NotEmpty(t, DatabaseDSN())
			c.DatabaseDSN = ""
		})
	}
}

func TestFileStoragePath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "FileStoragePath can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithFileStoragePath("test"))
			assert.NotEmpty(t, FileStoragePath())
			c.FileStoragePath = ""
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		want *config
		name string
	}{
		{
			name: "Pointer to new config will be returned",
			want: &config{
				ServerAddress:   "localhost:8080",
				BaseURL:         "http://localhost:8080",
				FileStoragePath: "",
				DatabaseDSN:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewConfig())
		})
	}
}

func TestServerAddress(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ServerAddress can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewConfig()
			assert.NotEmpty(t, ServerAddress())
		})
	}
}

func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config WithBaseURL can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithBaseURL("test"))
			assert.Equal(t, c.BaseURL, "test")
			c.BaseURL = "http://localhost:8080"
		})
	}
}

func TestWithDatabaseConnection(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config WithDatabaseConnection can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithDatabaseConnection("test"))
			assert.Equal(t, c.DatabaseDSN, "test")
			c.DatabaseDSN = ""
		})
	}
}

func TestWithFileStoragePath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config WithFileStoragePath can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithFileStoragePath("test"))
			assert.Equal(t, c.FileStoragePath, "test")
			c.FileStoragePath = ""
		})
	}
}

func TestWithServerAddress(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config WithServerAddress can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithServerAddress("test"))
			assert.Equal(t, c.ServerAddress, "test")
			c.ServerAddress = "localhost:8080"
		})
	}
}
