package service

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type expandStorageMock struct {
	IsKeyFoundInStore bool
}

func (sm *expandStorageMock) Stats() (int, int, error) {
	return 0, 0, nil
}

func (sm *expandStorageMock) LinksByUUID(uuid string) ([]storage.UserURLs, bool) {
	if sm.IsKeyFoundInStore {
		var sl []storage.UserURLs

		sl = append(sl, storage.UserURLs{
			ShortURL:    "test",
			OriginalURL: "https://github.com/",
		})

		return sl, true
	}
	return nil, false
}

func (sm *expandStorageMock) Store(key *string, url string, uid string) {}
func (sm *expandStorageMock) Get(key string) (string, bool, bool) {
	expandedURL := "https://github.com/"
	if !sm.IsKeyFoundInStore {
		expandedURL = ""
	}
	return expandedURL, sm.IsKeyFoundInStore, false
}

func TestNewURLExpandService(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		args args
		want *URLExpandService
		name string
	}{
		{
			name: "URLURLExpandService can be created",
			args: args{
				storage: &expandStorageMock{},
			},
			want: &URLExpandService{
				storage: &expandStorageMock{},
				logger:  &zap.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewURLExpandService(tt.args.storage, &zap.Logger{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewURLExpandService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLExpandService_ExpandURL(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Service can retrieve expanded URL by it's key",
			fields:  fields{storage: &expandStorageMock{IsKeyFoundInStore: true}},
			args:    args{key: "key"},
			want:    "https://github.com/",
			wantErr: false,
		},
		{
			name:    "Service will throw error on retrieve expanded URL by it's key if key does not exists in storage",
			fields:  fields{storage: &expandStorageMock{IsKeyFoundInStore: false}},
			args:    args{key: "key"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URLExpandService{
				storage: tt.fields.storage,
			}
			got, err := u.ExpandURL(tt.args.key)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestURLExpandService_ExpandUserLinks(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Service can retrieve expanded URL by it's uuid",
			fields:  fields{storage: &expandStorageMock{IsKeyFoundInStore: true}},
			args:    args{key: "1"},
			want:    "https://github.com/",
			wantErr: false,
		},
		{
			name:    "Service will throw error on retrieve expanded URL by it's uuid if key does not exists in storage",
			fields:  fields{storage: &expandStorageMock{IsKeyFoundInStore: false}},
			args:    args{key: "1"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URLExpandService{
				storage: tt.fields.storage,
			}
			got, err := u.ExpandUserLinks(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, got)
				assert.NoError(t, err)
			}
		})
	}
}
