package handlers

import (
	"net/http"

	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
)

type URLDeleteHandler struct {
	service service.URLDelete
}

func NewURLDeleteHandler(service service.URLDelete) *URLDeleteHandler {
	return &URLDeleteHandler{
		service: service,
	}
}

func (h *URLDeleteHandler) Delete(w http.ResponseWriter, req *http.Request) {
	err := h.service.Delete(req)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
