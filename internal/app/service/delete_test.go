package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type DBMock struct {
	hasError  bool
	hasConn   bool
	isConnNil bool
}

func (d *DBMock) Stats() (int, int, error) {
	return 0, 0, nil
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

func TestNewURLDeleteService(t *testing.T) {
	type args struct {
		storage storage.DB
		l       *zap.Logger
	}
	tests := []struct {
		args args
		want *URLDeleteService
		name string
	}{
		{
			name: "UrlDeleteService can be created",
			args: args{
				storage: &DBMock{},
				l:       &zap.Logger{},
			},
			want: &URLDeleteService{
				storage: &DBMock{},
				logger:  &zap.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLDeleteService(tt.args.storage, tt.args.l), "NewURLDeleteService(%v, %v)", tt.args.storage, tt.args.l)
		})
	}
}

func Test_generateCh(t *testing.T) {
	type args struct {
		uid  string
		data []string
	}
	tests := []struct {
		want chan storage.BatchDelete
		name string
		args args
	}{
		{
			name: "Can generate Ch",
			args: args{
				uid:  "1",
				data: []string{"1", "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, generateCh(tt.args.uid, tt.args.data))
		})
	}
}
