package service

import (
	"errors"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
)

var _ interfaces.URLExpand = (*URLExpandService)(nil)

type URLExpandService struct {
	storage interfaces.Storage
}

func NewURLExpandService(storage interfaces.Storage) *URLExpandService {
	return &URLExpandService{
		storage: storage,
	}
}

func (u *URLExpandService) ExpandURL(key string) (string, error) {
	url, ok := u.storage.Get(key)
	if !ok || url == "" {
		return url, errors.New("error in expanding shortened link")
	}

	return url, nil
}
