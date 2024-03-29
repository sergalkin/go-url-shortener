package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type URLExpandHandlerMock struct {
	hasErrorInExpandingURL bool
}

func (u *URLExpandHandlerMock) ExpandUserLinks(uuid string) ([]storage.UserURLs, error) {
	return nil, nil
}

func (u *URLExpandHandlerMock) ExpandURL(key string) (string, error) {
	if u.hasErrorInExpandingURL {
		return "", errors.New("error")
	}
	return "https://yandex.ru", nil
}

func TestNewURLExpandHandler(t *testing.T) {
	type args struct {
		service service.URLExpand
	}
	tests := []struct {
		args args
		want *URLExpandHandler
		name string
	}{
		{
			name: "URLExpandHandler can be created",
			args: args{
				service: &URLExpandHandlerMock{},
			},
			want: &URLExpandHandler{
				service: &URLExpandHandlerMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLExpandHandler(tt.args.service), "NewURLShortenerHandler(%v)", tt.args.service)
		})
	}
}

func TestURLExpandHandler_ExpandURL(t *testing.T) {
	type want struct {
		contentType string
		response    string
		code        int
	}

	tests := []struct {
		name       string
		body       string
		urlHandler *URLExpandHandlerMock
		want       want
	}{
		{
			name: "On making GET request with proper short URL user will be redirect to long URL",
			body: "key",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			urlHandler: &URLExpandHandlerMock{
				hasErrorInExpandingURL: false,
			},
		},
		{
			name: "On making GET request with non existing short URL server will respond with status code 404 and error message",
			body: "key",
			want: want{
				code:        http.StatusNotFound,
				response:    "error\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlHandler: &URLExpandHandlerMock{
				hasErrorInExpandingURL: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/{id}", NewURLExpandHandler(tt.urlHandler).ExpandURL)

			ts := httptest.NewServer(r)

			resp, body := expandTestRequest(t, ts, http.MethodGet, "/"+tt.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)

			if tt.want.code == http.StatusNotFound {
				assert.Equal(t, tt.want.response, body)
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			} else {
				assert.Equal(t, "yandex.ru", resp.Request.URL.Host)
			}
		})
	}
}

func TestURLExpandHandler_UserURLs(t *testing.T) {
	type fields struct {
		service service.URLExpand
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "can retrieve list of URLs",
			fields: fields{service: &URLExpandHandlerMock{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(middleware.Cookie)
			r.Get("/user/urls", NewURLExpandHandler(tt.fields.service).UserURLs)

			ts := httptest.NewServer(r)
			resp, body := expandTestRequest(t, ts, http.MethodGet, "/user/urls")
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func expandTestRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
