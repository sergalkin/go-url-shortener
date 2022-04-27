package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type DBHandler struct {
	storage storage.DB
}

func NewDBHandler(storage storage.DB) *DBHandler {
	return &DBHandler{
		storage: storage,
	}
}

func (h *DBHandler) Ping(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	err := h.storage.Ping(ctx)
	if err != nil {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
