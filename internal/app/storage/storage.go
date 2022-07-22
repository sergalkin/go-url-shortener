package storage

import (
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
)

type Storage interface {
	Store(key *string, url string)
	Get(key string) (string, bool, bool)
	LinksByUUID(uuid string) ([]UserURLs, bool)
}

func NewStorage(l *zap.Logger) (Storage, error) {
	switch {
	case config.DatabaseDSN() != "":
		return NewDBConnection(l, true)
	case config.FileStoragePath() != "":
		return NewFile(config.FileStoragePath(), l), nil
	default:
		return NewMemory(l), nil
	}
}
