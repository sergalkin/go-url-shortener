package shortener

import (
	"errors"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type URLServiceMock struct {
	hasErrorInShortenURL   bool
	hasErrorInExpandingURL bool
}

func (u *URLServiceMock) ShortenURL(url string) (string, error) {
	if u.hasErrorInShortenURL {
		return "", errors.New("error")
	}
	return "randomKey", nil
}
func (u *URLServiceMock) ExpandURL(key string) (string, error) {
	if u.hasErrorInExpandingURL {
		return "", errors.New("error")
	}
	return "https://yandex.ru", nil
}

func TestNewURLShortenerHandler(t *testing.T) {
	type args struct {
		service interfaces.URLService
	}
	tests := []struct {
		name string
		args args
		want *URLShortenerHandler
	}{
		{
			name: "URLShortenHandler can be created",
			args: args{
				service: &URLServiceMock{},
			},
			want: &URLShortenerHandler{
				&URLServiceMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLShortenerHandler(tt.args.service), "NewURLShortenerHandler(%v)", tt.args.service)
		})
	}
}

func TestURLShortenerHandler_URLHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		body          string
		want          want
		urlService    *URLServiceMock
		requestMethod string
	}{
		{
			name: "HEAD method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodHead,
		},
		{
			name: "PUT method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodPut,
		},
		{
			name: "PATCH method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodPatch,
		},
		{
			name: "DELETE method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodDelete,
		},
		{
			name: "TRACE method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodTrace,
		},
		{
			name: "OPTIONS method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodOptions,
		},
		{
			name: "CONNECT method not allowed",
			body: "",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodConnect,
		},
		{
			name: "ShortenURL will be called on POST request",
			body: "",
			want: want{
				code:        http.StatusUnprocessableEntity,
				response:    "Body must have a link\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   true,
				hasErrorInExpandingURL: false,
			},
			requestMethod: http.MethodPost,
		},
		{
			name: "ExpandURL will be called on GET request",
			body: "",
			want: want{
				code:        http.StatusNotFound,
				response:    "error\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: true,
			},
			requestMethod: http.MethodGet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.requestMethod, "/", body)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(NewURLShortenerHandler(tt.urlService).URLHandler)
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, w.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestURLShortenerHandler_expandURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name       string
		body       string
		want       want
		urlService *URLServiceMock
	}{
		{
			name: "On making GET request with proper short URL user will be redirect to long URL",
			body: "",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
			},
		},
		{
			name: "On making GET request with non existing short URL server will respond with status code 404 and error message",
			body: "",
			want: want{
				code:        http.StatusNotFound,
				response:    "error\n",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodGet, "/", body)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(NewURLShortenerHandler(tt.urlService).expandURL)
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, w.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestURLShortenerHandler_shortenURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name       string
		body       string
		want       want
		urlService *URLServiceMock
	}{
		{
			name: "On making POST request with proper body service will generate short URL and return it in response",
			body: "https://yandex.ru",
			want: want{
				code:        http.StatusCreated,
				response:    "http://localhost:8080/randomKey",
				contentType: "text/plain; charset=utf-8",
			},
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
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
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   false,
				hasErrorInExpandingURL: false,
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
			urlService: &URLServiceMock{
				hasErrorInShortenURL:   true,
				hasErrorInExpandingURL: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(NewURLShortenerHandler(tt.urlService).shortenURL)
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, w.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
