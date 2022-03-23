package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
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
			got, ok := m.Get(tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestMemory_Store(t *testing.T) {
	type fields struct {
		urls map[string]string
	}
	type args struct {
		key string
		url string
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		expectedLength   int
		expectedElements fields
	}{
		{
			name: "Long URL can be stored in memory struct by it's short encoded sequence",
			fields: fields{
				urls: map[string]string{},
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
			m := &Memory{
				urls: tt.fields.urls,
			}
			m.Store(tt.args.key, tt.args.url)

			assert.Len(t, m.urls, tt.expectedLength)
			reflect.DeepEqual(m.urls, tt.expectedElements)
		})
	}
}

func TestNewMemory(t *testing.T) {
	tests := []struct {
		name string
		want *Memory
	}{
		{
			name: "Memory object can be created",
			want: &Memory{urls: map[string]string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewMemory())
		})
	}
}
