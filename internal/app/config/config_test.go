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

func TestWithEnableHTTPS(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config WithEnableHTTPS can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithEnableHTTPS(true))
			assert.True(t, EnableHTTPS())
			c.EnableHTTPS = false
		})
	}
}

func TestEnableHTTPS(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "EnableHTTPS can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig()
			c.EnableHTTPS = true
			assert.True(t, EnableHTTPS())
			c.EnableHTTPS = false
		})
	}
}

func TestJSONConfigPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "JSONConfigPath can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig()
			c.JSONConfigPath = "test"
			assert.Equal(t, "test", JSONConfigPath())
			c.JSONConfigPath = ""
		})
	}
}

func TestWithJSONConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config WithEnableHTTPS can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithJSONConfig("test"))
			assert.Equal(t, "test", JSONConfigPath())
			c.JSONConfigPath = ""
		})
	}
}

func Test_config_SetJSONValues(t *testing.T) {
	type fields struct {
		ServerAddress   string
		BaseURL         string
		FileStoragePath string
		DatabaseDSN     string
		JSONConfigPath  string
		TrustedSubnet   string
		EnableHTTPS     bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Can set values from json file",
			fields: fields{
				JSONConfigPath: "config.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config{
				JSONConfigPath: tt.fields.JSONConfigPath,
			}
			assert.False(t, c.EnableHTTPS)
			c.SetJSONValues()
			assert.True(t, c.EnableHTTPS)
		})
	}
}

func TestWithTrustedSubnet(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config TrustedSubnet can be created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(WithTrustedSubnet("test"))
			assert.Equal(t, "test", TrustedSubnet())
			c.TrustedSubnet = ""
		})
	}
}

func TestTrustedSubnet(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TrustedSubnet can be retrieved from config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig()
			c.TrustedSubnet = "test"
			assert.Equal(t, "test", TrustedSubnet())
			c.TrustedSubnet = ""
		})
	}
}
