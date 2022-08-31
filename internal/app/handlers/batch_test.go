package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type BatchHandlerMock struct {
}

func TestNewBatchHandler(t *testing.T) {
	type args struct {
		storage storage.DB
		l       *zap.Logger
	}
	tests := []struct {
		args args
		want *BatchHandler
		name string
	}{
		{
			name: "DBHandler can be created",
			args: args{
				storage: &DBMock{},
				l:       zap.NewNop(),
			},
			want: &BatchHandler{
				storage: &DBMock{},
				logger:  zap.NewNop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewBatchHandler(tt.args.storage, tt.args.l), "NewBatchHandler(%v, %v)", tt.args.storage, tt.args.l)
		})
	}
}

func TestBatchHandler_BatchInsert_ThrowErr(t *testing.T) {
	tests := []struct {
		name    string
		handler BatchHandler
	}{
		{
			name: "batch insert will throw error if cant get uid from cookie",
			handler: BatchHandler{
				storage: &DBMock{},
				logger:  zap.NewNop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/api/shorten/batch", tt.handler.BatchInsert)

			ts := httptest.NewServer(r)

			resp, _ := batchTestRequest(t, ts, http.MethodPost, "/api/shorten/batch",
				strings.NewReader(`[{"correlation_id":"66f29390-381c-4a6a-9df9-74a0247ebe72", "original_url": "test.ya.ru"}]`),
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})
	}
}

func TestBatchHandler_BatchInsert(t *testing.T) {
	tests := []struct {
		name    string
		handler BatchHandler
	}{
		{
			name: "can batch insert",
			handler: BatchHandler{
				storage: &DBMock{},
				logger:  zap.NewNop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(middleware.Cookie)
			r.Post("/api/shorten/batch", tt.handler.BatchInsert)

			ts := httptest.NewServer(r)

			resp, body := batchTestRequest(t, ts, http.MethodPost, "/api/shorten/batch",
				strings.NewReader(`[{"correlation_id":"66f29390-381c-4a6a-9df9-74a0247ebe72", "original_url": "test.ya.ru"}]`),
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			assert.NotEmpty(t, body)
		})
	}
}

func batchTestRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
