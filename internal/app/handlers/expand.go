package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type URLExpandHandler struct {
	service service.URLExpand
}

func NewURLExpandHandler(service service.URLExpand) *URLExpandHandler {
	return &URLExpandHandler{
		service: service,
	}
}

func (h *URLExpandHandler) ExpandURL(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")

	originalLink, err := h.service.ExpandURL(key)

	if errors.Is(err, utils.ErrLinkIsDeleted) {
		http.Error(w, err.Error(), http.StatusGone)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLExpandHandler) UserURLs(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var uuid string
	err := utils.Decode(middleware.GetUUID(), &uuid)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	links, errExpand := h.service.ExpandUserLinks(uuid)
	if errExpand != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for k, v := range links {
		ss := strings.Split(v.ShortURL, "/")
		links[k].ShortURL = config.BaseURL() + "/" + ss[len(ss)-1]
	}

	w.WriteHeader(http.StatusOK)
	errEnc := json.NewEncoder(w).Encode(links)
	if errEnc != nil {
		utils.JSONError(w, errEnc.Error(), http.StatusInternalServerError)
		return
	}
}
