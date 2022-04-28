package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
)

func TestNewFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want *fileStore
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

	err := os.Remove("tmp")
	if err != nil {
		log.Fatalln(err)
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
			got, ok := m.Get(tt.args.key)
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
			fs := NewFile(tt.fields.filePath, &zap.Logger{})
			fs.urls = tt.fields.urls

			fs.Store(&tt.args.key, tt.args.url)
			assert.NotEmpty(t, fs.urls)
		})
	}

	err := os.Remove("tmp")
	if err != nil {
		log.Fatalln(err)
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

	err := os.Remove("tmp")
	if err != nil {
		log.Fatalln(err)
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

	err := os.Remove("tmp")
	if err != nil {
		log.Fatalln(err)
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
		marshal, err := json.Marshal(urlRecord{Key: k, URL: v})
		if err != nil {
			log.Fatalln(err)
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
