package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type URLDelete interface {
	Delete(r *http.Request) error
}

var _ URLDelete = (*URLDeleteService)(nil)

type URLDeleteService struct {
	storage storage.DB
	logger  *zap.Logger
}

func NewURLDeleteService(storage storage.DB, l *zap.Logger) *URLDeleteService {
	return &URLDeleteService{
		storage: storage,
		logger:  l,
	}
}

func (s *URLDeleteService) Delete(r *http.Request) error {
	uid, data, err := getDataFromBody(s, r)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		ch := generateCh(uid, data)
		s.storage.DeleteThroughCh(ch)
	}

	return nil
}

func getDataFromBody(s *URLDeleteService, r *http.Request) (string, []string, error) {
	var uid string
	err := utils.Decode(middleware.GetUUID(), &uid)
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return "", nil, err
	}

	b, errB := ioutil.ReadAll(r.Body)
	if errB != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return "", nil, err
	}

	var arr []string
	err = json.Unmarshal(b, &arr)
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return "", nil, err
	}

	return uid, arr, nil
}

func generateCh(uid string, data []string) chan storage.BatchDelete {
	inputCh := make(chan storage.BatchDelete)

	go func() {
		inputCh <- storage.BatchDelete{
			UID: uid,
			Arr: data,
		}
	}()

	return inputCh
}
