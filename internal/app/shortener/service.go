package shortener

import (
	"errors"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	sequence "github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

var _ interfaces.URLService = (*URLShortenerService)(nil)

type URLShortenerService struct {
	storage interfaces.Storage
}

func NewURLShortenerService(storage interfaces.Storage) *URLShortenerService {
	return &URLShortenerService{
		storage: storage,
	}
}

func (u *URLShortenerService) ShortenURL(url string) string {
	for {
		key := sequence.Generate(8)

		_, ok := u.storage.Get(key)
		if !ok {
			u.storage.Store(key, url)
			return key
		}
	}
}

func (u *URLShortenerService) ExpandURL(key string) (string, error) {
	url, ok := u.storage.Get(key)
	if !ok || url == "" {
		return url, errors.New("error in expanding shortened link")
	}

	return url, nil
}
