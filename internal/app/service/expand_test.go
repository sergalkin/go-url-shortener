package service

import (
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type expandStorageMock struct {
	IsKeyFoundInStore bool
}

func (sm *expandStorageMock) Store(key string, url string) {}
func (sm *expandStorageMock) Get(key string) (string, bool) {
	expandedURL := "https://github.com/"
	if !sm.IsKeyFoundInStore {
		expandedURL = ""
	}
	return expandedURL, sm.IsKeyFoundInStore
}

func TestNewURLExpandService(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		name string
		args args
		want *URLExpandService
	}{
		{
			name: "URLURLExpandService can be created",
			args: args{
				storage: &expandStorageMock{},
			},
			want: &URLExpandService{
				storage: &expandStorageMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewURLExpandService(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
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
