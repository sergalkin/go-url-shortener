package service

import (
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

var _ Internal = (*InternalService)(nil)

type Internal interface {
	Stats() (int, int, error)
}

type InternalService struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewInternalService(storage storage.Storage, l *zap.Logger) *InternalService {
	return &InternalService{
		storage: storage,
		logger:  l,
	}
}

// Stats - returns users and urls via proper storage manager.
func (i *InternalService) Stats() (int, int, error) {
	return i.storage.Stats()
}
