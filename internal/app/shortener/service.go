package shortener

import (
	"errors"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
)

var _ interfaces.URLService = (*URLShortenerService)(nil)

type URLShortenerService struct {
	storage interfaces.Storage
	seq     interfaces.Sequence
}

func NewURLShortenerService(storage interfaces.Storage, seq interfaces.Sequence) *URLShortenerService {
	return &URLShortenerService{
		storage: storage,
		seq:     seq,
	}
}

func (u *URLShortenerService) ShortenURL(url string) (string, error) {
	for {
		key, err := u.seq.Generate(8)
		if err != nil {
			return "", err
		}

		_, ok := u.storage.Get(key)
		if !ok {
			u.storage.Store(key, url)
			return key, nil
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
