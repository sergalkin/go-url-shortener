package service

import (
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

var _ URLShorten = (*URLShortenerService)(nil)

type URLShorten interface {
	ShortenURL(url string) (string, error)
}

type URLShortenerService struct {
	storage storage.Storage
	seq     utils.SequenceGenerator
	logger  *zap.Logger
}

func NewURLShortenerService(storage storage.Storage, seq utils.SequenceGenerator, l *zap.Logger) *URLShortenerService {
	return &URLShortenerService{
		storage: storage,
		seq:     seq,
		logger:  l,
	}
}

func (u *URLShortenerService) ShortenURL(url string) (string, error) {
	for {
		key, err := u.seq.Generate(8)
		if err != nil {
			u.logger.Error(err.Error(), zap.Error(err))
			return "", err
		}

		_, ok := u.storage.Get(key)
		if !ok {
			keyBeforeStore := key
			u.storage.Store(&key, url)

			if keyBeforeStore != key {
				return key, utils.ErrLinksConflict
			}

			return key, nil
		}
	}
}
