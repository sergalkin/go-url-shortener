package handlers

import (
	"encoding/json"
	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
	"io"
	"net/http"
)

type URLShortenerHandler struct {
	service service.URLShorten
}

func NewURLShortenerHandler(service service.URLShorten) *URLShortenerHandler {
	return &URLShortenerHandler{
		service: service,
	}
}

func (h *URLShortenerHandler) ShortenURL(w http.ResponseWriter, req *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(req.Body)
	bodyReq, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := string(bodyReq)
	if len(body) == 0 {
		http.Error(w, "Body must have a link", http.StatusUnprocessableEntity)
		return
	}

	key, shortenErr := h.service.ShortenURL(body)
	if shortenErr != nil {
		http.Error(w, shortenErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.BaseURL() + key))
}

func (h *URLShortenerHandler) APIShortenURL(w http.ResponseWriter, req *http.Request) {
	requestData := struct {
		URL string
	}{}

	responseData := struct {
		Result string `json:"result,omitempty"`
	}{}

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.URL == "" {
		utils.JSONError(w, "Body must have a link", http.StatusUnprocessableEntity)
		return
	}

	key, shortenErr := h.service.ShortenURL(requestData.URL)

	if shortenErr != nil {
		utils.JSONError(w, shortenErr.Error(), http.StatusInternalServerError)
		return
	}

	responseData.Result = config.BaseURL() + key

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseData)
}
