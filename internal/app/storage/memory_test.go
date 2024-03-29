package storage

import (
	"reflect"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemory_Get(t *testing.T) {
	type fields struct {
		urls map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		ok     bool
	}{
		{
			name: "Long URL can be retrieved by it's short encoded sequence",
			fields: fields{
				urls: map[string]string{"randomKey": "https://yandex.ru/"},
			},
			args: args{key: "randomKey"},
			want: "https://yandex.ru/",
			ok:   true,
		},
		{
			name: "Empty string and false status will be returned on accessing empty map",
			fields: fields{
				urls: map[string]string{},
			},
			args: args{key: "randomKey"},
			want: "",
			ok:   false,
		},
		{
			name: "Empty string and false status will be returned on accessing non existing element",
			fields: fields{
				urls: map[string]string{"key": "https://yandex.ru/"},
			},
			args: args{key: "randomKey"},
			want: "",
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				urls: tt.fields.urls,
			}
			got, ok, _ := m.Get(tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestMemory_Store(t *testing.T) {
	type fields struct {
		urls     map[string]string
		userURLs map[string][]UserURLs
	}
	type args struct {
		key string
		url string
	}
	tests := []struct {
		expectedElements fields
		fields           fields
		args             args
		name             string
		expectedLength   int
	}{
		{
			name: "Long URL can be stored in memory struct by it's short encoded sequence",
			fields: fields{
				urls:     map[string]string{},
				userURLs: map[string][]UserURLs{},
			},
			args:             args{"randomKey", "https://yandex.ru/"},
			expectedLength:   1,
			expectedElements: fields{urls: map[string]string{"randomKey": "https://yandex.ru/"}},
		},
		{
			name: "Long URL can be stored in memory struct by it's short encoded sequence even if URLs map not empty",
			fields: fields{
				urls: map[string]string{
					"firstKey": "https://test.ru/test",
				},
				userURLs: map[string][]UserURLs{},
			},
			args:           args{"randomKey", "https://yandex.ru/"},
			expectedLength: 2,
			expectedElements: fields{urls: map[string]string{
				"randomKey": "https://yandex.ru/",
				"firstKey":  "https://test.ru/test",
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(zap.NewNop())
			m.urls = tt.fields.urls

			m.Store(&tt.args.key, tt.args.url, "046cf584-df95-43fd-a2fc-f95a85c7bb95")

			assert.Len(t, m.urls, tt.expectedLength)
			reflect.DeepEqual(m.urls, tt.expectedElements)
		})
	}
}

func TestNewMemory(t *testing.T) {
	tests := []struct {
		want *Memory
		name string
	}{
		{
			name: "Memory object can be created",
			want: &Memory{
				urls:     map[string]string{},
				userURLs: map[string][]UserURLs{},
				logger:   &zap.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewMemory(&zap.Logger{}))
		})
	}
}

func TestMemory_LinksByUUID(t *testing.T) {
	type fields struct {
		urls     map[string]string
		userURLs map[string][]UserURLs
		logger   *zap.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   []UserURLs
	}{
		{
			name: "Can retrieve links UUID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				urls:     tt.fields.urls,
				userURLs: map[string][]UserURLs{},
			}

			m.userURLs["1"] = append(m.userURLs["1"], UserURLs{ShortURL: "", OriginalURL: ""})
			links, ok := m.LinksByUUID("1")

			assert.True(t, ok)
			assert.NotEmpty(t, links)
		})
	}
}

func TestMemory_Stats(t *testing.T) {
	type fields struct {
		logger   *zap.Logger
		urls     map[string]string
		userURLs map[string][]UserURLs
	}
	tests := []struct {
		fields fields
		name   string
	}{
		{
			name: "Stats can be retrieved via file in memory manager",
			fields: fields{
				logger: zap.NewNop(),
				urls:   map[string]string{"test": "ya.ru", "test2": "vk.ru"},
				userURLs: map[string][]UserURLs{"test": append(make([]UserURLs, 1), UserURLs{
					ShortURL:    "short",
					OriginalURL: "long",
				}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				logger:   tt.fields.logger,
				urls:     tt.fields.urls,
				userURLs: tt.fields.userURLs,
			}

			got, got1, err := m.Stats()
			assert.Equal(t, 2, got)
			assert.Equal(t, 1, got1)
			assert.NoError(t, err)
		})
	}
}
