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

	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
)

type InternalHandlerMock struct {
	hasError bool
}

func (i *InternalHandlerMock) Stats() (int, int, error) {
	if i.hasError {
		return 0, 0, errors.New("error")
	}
	return 1, 2, nil
}

func TestInternalHandler_Stats(t *testing.T) {
	type want struct {
		contentType string
		response    string
		code        int
	}

	tests := []struct {
		name       string
		body       string
		urlHandler *InternalHandlerMock
		want       want
	}{
		{
			name: "On making GET request will retrieve stats if has no errors.",
			want: want{
				code:        http.StatusOK,
				response:    "{\"urls\":1,\"users\":2}\n",
				contentType: "application/json; charset=utf-8",
			},
			urlHandler: &InternalHandlerMock{
				hasError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/api/internal/stats", NewInternalHandler(tt.urlHandler).Stats)

			ts := httptest.NewServer(r)

			resp, body := internalTestRequest(t, ts, http.MethodGet, "/api/internal/stats")
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, body)
		})
	}
}

func TestNewInternalHandler(t *testing.T) {
	type args struct {
		s service.Internal
	}
	tests := []struct {
		args args
		want *InternalHandler
		name string
	}{
		{
			name: "InternalHandler can be created",
			args: args{
				s: &InternalHandlerMock{},
			},
			want: &InternalHandler{
				service: &InternalHandlerMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInternalHandler(tt.args.s), "NewInternalHandler(%v)", tt.args.s)
		})
	}
}

func internalTestRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	req.Header.Add("X-Real-IP", "127.0.0.1")
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
