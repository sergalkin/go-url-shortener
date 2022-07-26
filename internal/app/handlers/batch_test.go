package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

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
