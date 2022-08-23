package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type DBMock struct {
	hasError  bool
	hasConn   bool
	isConnNil bool
}

func (d *DBMock) Ping(ctx context.Context) error {
	if d.hasError {
		return errors.New("error")
	}

	return nil
}

func (d *DBMock) Close(ctx context.Context) error {
	return nil
}

func (d *DBMock) Store(key *string, url string, uid string) {
}
func (d *DBMock) Get(key string) (string, bool, bool) {
	return "", true, true
}
func (d *DBMock) LinksByUUID(uuid string) ([]storage.UserURLs, bool) {
	return nil, false
}
func (d *DBMock) BatchInsert([]storage.BatchRequest) ([]storage.BatchLink, error) {
	return nil, nil
}
func (d *DBMock) SoftDeleteUserURLs(uuid string, ids []string) error {
	return nil
}

func (d *DBMock) DeleteThroughCh(channels ...chan storage.BatchDelete) {
}

func (d *DBMock) HasNotNilConn() bool {
	return d.isConnNil
}

func TestDBHandler_Ping(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		urlHandler *DBMock
		name       string
		body       string
		want       want
	}{
		{
			name: "Can return 200 status",
			want: want{code: http.StatusOK},
			urlHandler: &DBMock{
				hasError: false,
			},
		},
		{
			name: "Can return 500 status",
			want: want{code: http.StatusInternalServerError},
			urlHandler: &DBMock{
				hasError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Delete("/api/user/urls", NewDBHandler(tt.urlHandler, zap.NewNop()).Ping)

			ts := httptest.NewServer(r)

			resp, _ := shortenTestRequest(t, ts, http.MethodDelete, "/api/user/urls", strings.NewReader(tt.body))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

func TestNewDBHandler(t *testing.T) {
	type args struct {
		storage storage.DB
		l       *zap.Logger
	}
	tests := []struct {
		want *DBHandler
		args args
		name string
	}{
		{
			name: "DBHandler can be created",
			args: args{
				storage: &DBMock{},
				l:       zap.NewNop(),
			},
			want: &DBHandler{
				storage: &DBMock{},
				logger:  zap.NewNop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDBHandler(tt.args.storage, tt.args.l), "NewDBHandler(%v, %v)", tt.args.storage, tt.args.l)
		})
	}
}
