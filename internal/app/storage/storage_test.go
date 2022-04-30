package storage

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"

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

	err := os.Remove("tmp")
	if err != nil {
		log.Fatalln(err)
	}
}
