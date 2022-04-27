package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
)

type DbHandler struct {
	storage storage.Db
}

func NewDbHandler(storage storage.Db) *DbHandler {
	return &DbHandler{
		storage: storage,
	}
}

func (h *DbHandler) Ping(w http.ResponseWriter, req *http.Request) {
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
