package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/pkg/sequence"
)

func Test_db_Ping(t *testing.T) {
	tests := []struct {
		name string
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
			}

			cfg.DatabaseDSN = ""
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
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func BenchmarkDb_Store(b *testing.B) {
	cfg := config.NewConfig(
		config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
	)

	conn, err := NewDBConnection(zap.NewNop(), true)

	if err == nil {
		seq := sequence.NewSequence()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// do not stop here timer cause store need generated seq of letters, and it too takes time
			key, _ := seq.Generate(4)

			conn.Store(&key, "test.com")
		}
		b.StopTimer()
	}

	cfg.DatabaseDSN = ""
}

func Test_fanIn(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"can fanIn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make([]chan BatchDelete, 1)
			assert.NotNil(t, fanIn(ch...))
		})
	}
}
