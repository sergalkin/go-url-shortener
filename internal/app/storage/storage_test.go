package storage

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want Storage
		do   func()
	}{
		{
			name: "Memory storage will be created if no filepath is provided",
			want: &Memory{
				urls:     map[string]string{},
				userURLs: map[string][]UserURLs{},
				logger:   &zap.Logger{},
			},
			do: func() {},
		},
		{
			name: "FileStorage will be created if filepath is provided",
			want: &fileStore{
				urls:     map[string]string{},
				filePath: "tmp",
				userURLs: map[string][]UserURLs{},
				logger:   &zap.Logger{},
			},
			do: func() {
				config.NewConfig(config.WithFileStoragePath("tmp"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()
			s, _ := NewStorage(&zap.Logger{})
			assert.Equal(t, tt.want, s)
		})
	}

	if _, fErr := os.Stat("tmp"); fErr == nil {
		err := os.Remove("tmp")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func TestNewStorageCanCreateDbStorage(t *testing.T) {
	tests := []struct {
		name string
		do   func()
	}{
		{
			name: "DB will be created if DSN is provided",
			do: func() {
				config.NewConfig(
					config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()
			s, err := NewStorage(&zap.Logger{})
			if err == nil {
				assert.NotNil(t, s)
				assert.IsType(t, &db{}, s)
				assert.NotEmpty(t, s)
			}
		})
	}
}
