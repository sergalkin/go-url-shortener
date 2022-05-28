package service

import (
	"errors"
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type shortenStorageMock struct {
	IsKeyFoundInStore bool
}

func (sm *shortenStorageMock) LinksByUUID(uuid string) ([]storage.UserURLs, bool) {
	return nil, true
}

func (sm *shortenStorageMock) Store(key *string, url string) {}
func (sm *shortenStorageMock) Get(key string) (string, bool, bool) {
	expandedURL := "https://github.com/"
	if !sm.IsKeyFoundInStore {
		expandedURL = ""
	}
	return expandedURL, sm.IsKeyFoundInStore, false
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
		storage storage.Storage
		seq     utils.SequenceGenerator
	}
	tests := []struct {
		name string
		args args
		want *URLShortenerService
	}{
		{
			name: "",
			args: args{
				storage: &shortenStorageMock{},
				seq:     &sequenceMock{},
			},
			want: &URLShortenerService{
				storage: &shortenStorageMock{},
				seq:     &sequenceMock{},
				logger:  zap.NewNop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLShortenerService(tt.args.storage, tt.args.seq, zap.NewNop()), "NewURLShortenerService(%v, %v)", tt.args.storage, tt.args.seq)
		})
	}
}

func TestURLShortenerService_ShortenURL(t *testing.T) {
	type fields struct {
		storage storage.Storage
		seq     utils.SequenceGenerator
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
				storage: &shortenStorageMock{IsKeyFoundInStore: false},
				seq:     &sequenceMock{HasErrorInGenerationSeq: false},
			},
			args: args{url: "https://github.com/"},
		},
		{
			name: "Empty URL can be shortened and stored",
			fields: fields{
				storage: &shortenStorageMock{IsKeyFoundInStore: false},
				seq:     &sequenceMock{HasErrorInGenerationSeq: false},
			},
			args: args{url: ""},
		},
		{
			name: "If random sequence generator fails an error will thrown",
			fields: fields{
				storage: &shortenStorageMock{IsKeyFoundInStore: false},
				seq:     &sequenceMock{HasErrorInGenerationSeq: true},
			},
			args: args{url: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLShortenerService(tt.fields.storage, tt.fields.seq, zap.NewNop())

			got, err := u.ShortenURL(tt.args.url)
			if err != nil {
				assert.Empty(t, got)
			} else {
				assert.NotEmpty(t, got)
			}
		})
	}
}
