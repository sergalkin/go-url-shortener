package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type InternalStorageMock struct{}

func (i *InternalStorageMock) Store(key *string, url string, uid string) {
}

func (i *InternalStorageMock) Get(key string) (string, bool, bool) {
	return "test", true, true
}

func (i *InternalStorageMock) LinksByUUID(uuid string) ([]storage.UserURLs, bool) {
	return nil, false
}

func (i *InternalStorageMock) Stats() (int, int, error) {
	return 1, 2, nil
}

func TestNewInternalService(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		args args
		want *InternalService
		name string
	}{
		{
			name: "New Internal Service can be created",
			args: args{
				storage: &InternalStorageMock{},
			},
			want: &InternalService{
				storage: &InternalStorageMock{},
				logger:  zap.NewNop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInternalService(tt.args.storage, zap.NewNop()), "NewInternalService(%v)", tt.args.storage)
		})
	}
}

func TestInternalService_Stats(t *testing.T) {
	type fields struct {
		storage storage.Storage
		logger  *zap.Logger
	}
	tests := []struct {
		fields  fields
		wantErr assert.ErrorAssertionFunc
		name    string
		urls    int
		users   int
	}{
		{
			name: "Stats can be retrieved",
			fields: fields{
				storage: &InternalStorageMock{},
				logger:  zap.NewNop(),
			},
			urls:  1,
			users: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InternalService{
				storage: tt.fields.storage,
				logger:  tt.fields.logger,
			}

			got, got1, err := i.Stats()

			assert.NoError(t, err)
			assert.Equalf(t, tt.urls, got, "Stats()")
			assert.Equalf(t, tt.users, got1, "Stats()")
		})
	}
}
