package storage

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewFile(t *testing.T) {
	tests := []struct {
		want *fileStore
		name string
		path string
	}{
		{
			name: "File store can be created",
			path: "tmp",
			want: &fileStore{
				urls:     map[string]string{},
				filePath: "tmp",
				userURLs: map[string][]UserURLs{},
				logger:   &zap.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewFile(tt.path, &zap.Logger{}), "NewFile(%v)", tt.path)
		})
	}

	if _, fErr := os.Stat("tmp"); fErr == nil {
		err := os.Remove("tmp")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func Test_fileStore_Get(t *testing.T) {
	type fields struct {
		urls     map[string]string
		filePath string
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
			m := &fileStore{
				urls: tt.fields.urls,
			}
			got, ok, _ := m.Get(tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_fileStore_Store(t *testing.T) {
	type fields struct {
		urls     map[string]string
		userURLs map[string][]UserURLs
		filePath string
	}
	type args struct {
		key string
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "URL can be stored in file storage struct",
			fields: fields{
				urls:     map[string]string{},
				userURLs: map[string][]UserURLs{},
				filePath: "tmp",
			},
			args: args{
				key: "test",
				url: "https://test.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFile(tt.fields.filePath, zap.NewNop())
			fs.urls = tt.fields.urls

			fs.Store(&tt.args.key, tt.args.url, "046cf584-df95-43fd-a2fc-f95a85c7bb95")
			assert.NotEmpty(t, fs.urls)
		})
	}

	if _, fErr := os.Stat("tmp"); fErr == nil {
		err := os.Remove("tmp")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func Test_fileStore_loadFromFile(t *testing.T) {
	type fields struct {
		urls     map[string]string
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Links can be loaded from file",
			fields: fields{
				urls:     map[string]string{},
				filePath: "tmp",
			},
		},
	}
	writeTestingDataToFile("tmp")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fileStore{
				urls:     tt.fields.urls,
				filePath: tt.fields.filePath,
			}

			err := m.loadFromFile()
			if err != nil {
				log.Fatalln(err)
			}
			assert.NotEmpty(t, m.urls)
		})
	}

	if _, fErr := os.Stat("tmp"); fErr == nil {
		err := os.Remove("tmp")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func Test_fileStore_saveToFile(t *testing.T) {
	type args struct {
		key string
		url string
	}
	tests := []struct {
		name string
		path string
		args args
	}{
		{
			name: "Record can be saved to File",
			path: "tmp",
			args: args{key: "test", url: "http://yandex.ru"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fileStore{filePath: tt.path}
			m.saveToFile(tt.args.key, tt.args.url)
			assert.FileExists(t, "tmp")

			f, err := os.OpenFile("tmp", os.O_RDONLY, 0644)
			if err != nil {
				log.Fatalln(err)
			}

			assert.NotEmpty(t, f)
			err = f.Close()
			if err != nil {
				log.Fatalln(err)
			}
		})
	}

	if _, fErr := os.Stat("tmp"); fErr == nil {
		err := os.Remove("tmp")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func writeTestingDataToFile(filepath string) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	m := make(map[string]string)
	m["test"] = "http://test.ru"
	m["test2"] = "http://test2.ru"
	for k, v := range m {
		marshal, errMarsh := json.Marshal(urlRecord{Key: k, URL: v})
		if errMarsh != nil {
			log.Fatalln(errMarsh)
		}

		_, err = file.Write(marshal)
		file.Write([]byte("\n"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = file.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func Test_fileStore_LinksByUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Links can be restored from UUID",
			args: args{uuid: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linksMap := map[string][]UserURLs{}

			linksMap["1"] = append(linksMap["1"], UserURLs{ShortURL: "test", OriginalURL: "ya.ru"})
			m := &fileStore{userURLs: linksMap}

			got, ok := m.LinksByUUID(tt.args.uuid)
			assert.True(t, ok)
			assert.NotEmpty(t, got)
		})
	}
}

func Test_fileStore_Stats(t *testing.T) {
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
			name: "Stats can be retrieved via file manager",
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
			m := &fileStore{
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
