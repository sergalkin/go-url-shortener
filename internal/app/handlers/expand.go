package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	"net/http"
)

type URLExpandHandler struct {
	service interfaces.URLExpand
}

func NewURLExpandHandler(service interfaces.URLExpand) *URLExpandHandler {
	return &URLExpandHandler{
		service: service,
	}
}

func (h *URLExpandHandler) ExpandURL(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")

	originalLink, err := h.service.ExpandURL(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
