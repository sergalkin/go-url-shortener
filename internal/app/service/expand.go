package service

import (
	"errors"

	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type URLExpand interface {
	ExpandURL(key string) (string, error)
	ExpandUserLinks(uuid string) ([]storage.UserURLs, error)
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

// ExpandURL - attempts to retrieve original URL by its shortened value.
func (u *URLExpandService) ExpandURL(key string) (string, error) {
	url, ok, isDeleted := u.storage.Get(key)

	if !ok || (url == "" && !isDeleted) {
		return url, errors.New("error in expanding shortened link")
	}

	if isDeleted {
		return url, utils.ErrLinkIsDeleted
	}

	return url, nil
}

// ExpandUserLinks - attempts to retrieve original URLs by provided uuid.
func (u *URLExpandService) ExpandUserLinks(uuid string) ([]storage.UserURLs, error) {
	links, ok := u.storage.LinksByUUID(uuid)
	if !ok {
		return links, errors.New("this UUID doesnt have links")
	}

	return links, nil
}
