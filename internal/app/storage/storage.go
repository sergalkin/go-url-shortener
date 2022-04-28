package storage

import "github.com/sergalkin/go-url-shortener.git/internal/app/config"

type Storage interface {
	Store(key *string, url string)
	Get(key string) (string, bool)
	LinksByUUID(uuid string) ([]UserURLs, bool)
}

func NewStorage() (Storage, error) {
	switch {
	case config.DatabaseDSN() != "":
		return NewDBConnection()
	case config.FileStoragePath() != "":
		return NewFile(config.FileStoragePath()), nil
	default:
		return NewMemory(), nil
	}
}
