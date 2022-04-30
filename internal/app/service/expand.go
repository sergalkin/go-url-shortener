package service

import (
	"errors"

	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type URLExpand interface {
	ExpandURL(key string) (string, error)
	ExpandUserLinks() ([]storage.UserURLs, error)
}

var _ URLExpand = (*URLExpandService)(nil)

type URLExpandService struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewURLExpandService(storage storage.Storage, l *zap.Logger) *URLExpandService {
	return &URLExpandService{
		storage: storage,
		logger:  l,
	}
}

func (u *URLExpandService) ExpandURL(key string) (string, error) {
	url, ok := u.storage.Get(key)
	if !ok || url == "" {
		return url, errors.New("error in expanding shortened link")
	}

	return url, nil
}

func (u *URLExpandService) ExpandUserLinks() ([]storage.UserURLs, error) {
	var uuid string
	err := utils.Decode(middleware.GetUUID(), &uuid)
	if err != nil {
		u.logger.Error(err.Error(), zap.Error(err))
	}

	links, ok := u.storage.LinksByUUID(uuid)
	if !ok {
		return links, errors.New("this UUID doesnt have links")
	}

	return links, nil
}
