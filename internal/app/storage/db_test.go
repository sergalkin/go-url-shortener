package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
)

func Test_db_Ping(t *testing.T) {
	tests := []struct {
		name string
		do   func()
	}{
		{
			name: "Can ping DB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				errPing := conn.Ping(context.Background())

				assert.NoError(t, errPing)

				cfg.DatabaseDSN = ""
			}
		})
	}
}

func Test_db_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Can close DB conn without errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				errPing := conn.Close(context.Background())

				assert.NoError(t, errPing)

				cfg.DatabaseDSN = ""
			}
		})
	}
}
