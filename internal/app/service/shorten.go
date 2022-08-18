package service

import (
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
	"github.com/sergalkin/go-url-shortener.git/pkg/sequence"
)

var _ URLShorten = (*URLShortenerService)(nil)

type URLShorten interface {
	ShortenURL(url string) (string, error)
}

type URLShortenerService struct {
	storage storage.Storage
	seq     sequence.Generator
	logger  *zap.Logger
}

func NewURLShortenerService(storage storage.Storage, seq sequence.Generator, l *zap.Logger) *URLShortenerService {
	return &URLShortenerService{
		storage: storage,
		seq:     seq,
		logger:  l,
	}
}

// ShortenURL - shortens provided URL and stores it in storage.
func (u *URLShortenerService) ShortenURL(url string) (string, error) {
	var uuid string
	errD := utils.Decode(middleware.GetUUID(), &uuid)
	if errD != nil {
		u.logger.Error(errD.Error(), zap.Error(errD))
	}

	for {
		key, err := u.seq.Generate(8)
		if err != nil {
			u.logger.Error(err.Error(), zap.Error(err))
			return "", err
		}

		_, ok, _ := u.storage.Get(key)
		if !ok {
			keyBeforeStore := key
			u.storage.Store(&key, url, uuid)

			if keyBeforeStore != key {
				return key, utils.ErrLinksConflict
			}

			return key, nil
		}
	}
}
