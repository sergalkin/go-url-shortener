package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type InternalHandler struct {
	service service.Internal
}

type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

func NewInternalHandler(s service.Internal) *InternalHandler {
	return &InternalHandler{
		service: s,
	}
}

// Stats - will return count stored urls and users. Works only via trusted subnet.
func (h *InternalHandler) Stats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	urls, users, err := h.service.Stats()
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errEnc := json.NewEncoder(w).Encode(Stats{URLs: urls, Users: users})
	if errEnc != nil {
		utils.JSONError(w, errEnc.Error(), http.StatusInternalServerError)
		return
	}
}
