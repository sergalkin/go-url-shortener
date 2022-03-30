package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type URLShortenHandlerMock struct {
	hasErrorInShortenURL bool
}

func (u *URLShortenHandlerMock) ShortenURL(url string) (string, error) {
	if u.hasErrorInShortenURL {
		return "", errors.New("error")
	}
	return "randomKey", nil
}

func TestNewURLShortenerHandler(t *testing.T) {
	type args struct {
		service service.URLShorten
	}
	tests := []struct {
		name string
		args args
		want *URLShortenerHandler
	}{
		{
			name: "URLShortenHandler can be created",
			args: args{
				service: &URLShortenHandlerMock{},
			},
			want: &URLShortenerHandler{
				service: &URLShortenHandlerMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLShortenerHandler(tt.args.service), "NewURLShortenerHandler(%v)", tt.args.service)
		})
	}
}

func TestURLShortenerHandler_ShortenURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name       string
		body       string
		want       want
		urlHandler *URLShortenHandlerMock
	}{
		{
			name: "On making POST request with proper body service will generate short URL and return it in response",
			body: "https://yandex.ru",
			want: want{
				code:        http.StatusCreated,
				response:    "http://localhost:8080/randomKey",
				contentType: "text/plain; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: false,
			},
		},
		{
			name: "On making POST request without body service will return error message in response and 422 status code",
			body: "",
			want: want{
				code:        http.StatusUnprocessableEntity,
				response:    "Body must have a link\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: false,
			},
		},
		{
			name: "On making POST request with proper body service will return 500 status code and error message if short ULR can't be generated",
			body: "https://yandex.ru",
			want: want{
				code:        http.StatusInternalServerError,
				response:    "error\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/", NewURLShortenerHandler(tt.urlHandler).ShortenURL)

			ts := httptest.NewServer(r)

			resp, body := shortenTestRequest(t, ts, http.MethodPost, "/", strings.NewReader(tt.body))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestURLShortenerHandler_ApiShortenURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name       string
		body       string
		want       want
		urlHandler *URLShortenHandlerMock
	}{
		{
			name: "On making POST request with json body service will generate short URL and return it in response",
			body: `{ "url": "https://yandex.ru" }`,
			want: want{
				code:        http.StatusCreated,
				response:    "{\"result\":\"http://localhost:8080/randomKey\"}\n",
				contentType: "application/json; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: false,
			},
		},
		{
			name: "On making POST request without json body service will return error message in response and 400 status code",
			body: "",
			want: want{
				code:        http.StatusBadRequest,
				response:    "\"EOF\"\n",
				contentType: "application/json; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: false,
			},
		},
		{
			name: "On making POST request with empty url in json body service will return error message in response and 400 status code",
			body: `{ "url": "" }`,
			want: want{
				code:        http.StatusUnprocessableEntity,
				response:    "\"Body must have a link\"\n",
				contentType: "application/json; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: false,
			},
		},
		{
			name: "On making POST request with proper json body service will return 500 status code and error message if short ULR can't be generated",
			body: `{ "url": "https://yandex.ru" }`,
			want: want{
				code:        http.StatusInternalServerError,
				response:    "\"error\"\n",
				contentType: "application/json; charset=utf-8",
			},
			urlHandler: &URLShortenHandlerMock{
				hasErrorInShortenURL: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/api/shorten", NewURLShortenerHandler(tt.urlHandler).ApiShortenURL)

			ts := httptest.NewServer(r)

			resp, body := shortenTestRequest(t, ts, http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func shortenTestRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
