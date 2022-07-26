package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
)

type URLDeleteHandlerMock struct {
	hasError bool
}

func (h *URLDeleteHandlerMock) Delete(r *http.Request) error {
	if h.hasError {
		return errors.New("error")
	}

	return nil
}

func TestURLDeleteHandler_Delete(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		urlHandler *URLDeleteHandlerMock
		name       string
		body       string
		want       want
	}{
		{
			name: "Can return 202 status",
			want: want{code: http.StatusAccepted},
			urlHandler: &URLDeleteHandlerMock{
				hasError: false,
			},
		},
		{
			name: "Can return 204 status",
			want: want{code: http.StatusNoContent},
			urlHandler: &URLDeleteHandlerMock{
				hasError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Delete("/api/user/urls", NewURLDeleteHandler(tt.urlHandler).Delete)

			ts := httptest.NewServer(r)

			resp, _ := shortenTestRequest(t, ts, http.MethodDelete, "/api/user/urls", strings.NewReader(tt.body))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

func TestNewURLDeleteHandler(t *testing.T) {
	type args struct {
		service service.URLDelete
	}
	tests := []struct {
		want *URLDeleteHandler
		args args
		name string
	}{
		{
			name: "URLDeleteHandler can be created",
			args: args{
				service: &URLDeleteHandlerMock{},
			},
			want: &URLDeleteHandler{
				service: &URLDeleteHandlerMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLDeleteHandler(tt.args.service), "NewURLDeleteHandler(%v)", tt.args.service)
		})
	}
}
