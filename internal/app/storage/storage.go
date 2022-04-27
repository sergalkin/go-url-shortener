package storage

import "github.com/sergalkin/go-url-shortener.git/internal/app/config"

type Storage interface {
	Store(key string, url string)
	Get(key string) (string, bool)
	LinksByUUID(uuid string) ([]UserURLs, bool)
}

func NewStorage() Storage {
	fileStoragePath := config.FileStoragePath()

	if fileStoragePath == "" {
		return NewMemory()
	}

	return NewFile(fileStoragePath)
}
