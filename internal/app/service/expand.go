package service

import (
	"errors"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type URLExpand interface {
	ExpandURL(key string) (string, error)
}

var _ URLExpand = (*URLExpandService)(nil)

type URLExpandService struct {
	storage storage.Storage
}

func NewURLExpandService(storage storage.Storage) *URLExpandService {
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
