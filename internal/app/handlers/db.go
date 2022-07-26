package handlers

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type DBHandler struct {
	storage storage.DB
	logger  *zap.Logger
}

// NewDBHandler - creates DBHandler.
func NewDBHandler(storage storage.DB, l *zap.Logger) *DBHandler {
	return &DBHandler{
		storage: storage,
		logger:  l,
	}
}

// Ping - returns result of checking database availability to user.
func (h *DBHandler) Ping(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	err := h.storage.Ping(ctx)
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
