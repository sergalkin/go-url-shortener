package storage

import (
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
)

type Storage interface {
	// Store - store given URL into storage with key as id.
	Store(key *string, url string)
	// Get - trying to retrieve a URL from storage by provided key. As first bool value - returns was retrieval a
	// success or not and as second bool value return is link still present in storage.
	Get(key string) (string, bool, bool)
	// LinksByUUID - trying to retrieve slice of UserURLs. On successful retrieval returns true as bool value and false of
	// failure.
	LinksByUUID(uuid string) ([]UserURLs, bool)
}

// NewStorage - creates Storage implementation based on config options.
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
