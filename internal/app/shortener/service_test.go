package shortener

import (
	"errors"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type storageMock struct {
	IsKeyFoundInStore bool
}

func (sm *storageMock) Store(key string, url string) {}
func (sm *storageMock) Get(key string) (string, bool) {
	expandedURL := "https://github.com/"
	if !sm.IsKeyFoundInStore {
		expandedURL = ""
	}
	return expandedURL, sm.IsKeyFoundInStore
}

type sequenceMock struct {
	HasErrorInGenerationSeq bool
}

func (s *sequenceMock) Generate(lettersNumber int) (string, error) {
	if s.HasErrorInGenerationSeq {
		return "", errors.New("to generate random sequence positive number of letters must be provided")
	}

	return "randomString", nil
}

func TestNewURLShortenerService(t *testing.T) {
	type args struct {
		storage interfaces.Storage
	}
	tests := []struct {
		name string
		args args
		want *URLShortenerService
	}{
		{
			name: "URLShortenerService can be created",
			args: args{storage: &storageMock{}},
			want: &URLShortenerService{&storageMock{}, &sequenceMock{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewURLShortenerService(&storageMock{}, &sequenceMock{}))
		})
	}
}

func TestURLShortenerService_ExpandURL(t *testing.T) {
	type fields struct {
		storage interfaces.Storage
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
			fields:  fields{storage: &storageMock{IsKeyFoundInStore: true}},
			args:    args{key: "key"},
			want:    "https://github.com/",
			wantErr: false,
		},
		{
			name:    "Service will throw error on retrieve expanded URL by it's key if key does not exists in storage",
			fields:  fields{storage: &storageMock{IsKeyFoundInStore: false}},
			args:    args{key: "key"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URLShortenerService{
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

func TestURLShortenerService_ShortenURL(t *testing.T) {
	type fields struct {
		storage interfaces.Storage
		seq     interfaces.Sequence
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "URL can be shortened and stored",
			fields: fields{
				storage: &storageMock{IsKeyFoundInStore: false},
				seq:     &sequenceMock{HasErrorInGenerationSeq: false},
			},
			args: args{"https://github.com/"},
		},
		{
			name: "Empty URL can be shortened and stored",
			fields: fields{
				storage: &storageMock{IsKeyFoundInStore: false},
				seq:     &sequenceMock{HasErrorInGenerationSeq: false},
			},
			args: args{""},
		},
		{
			name: "If random sequence generator fails an error will thrown",
			fields: fields{
				storage: &storageMock{IsKeyFoundInStore: false},
				seq:     &sequenceMock{HasErrorInGenerationSeq: true},
			},
			args: args{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URLShortenerService{
				storage: tt.fields.storage,
				seq:     tt.fields.seq,
			}
			got, err := u.ShortenURL(tt.args.url)
			if err != nil {
				assert.Empty(t, got)
			} else {
				assert.NotEmpty(t, got)
			}
		})
	}
}
