package handlers

import (
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	"io"
	"net/http"
)

const host = "http://localhost:8080/"

type URLShortenerHandler struct {
	service interfaces.URLShorten
}

func NewURLShortenerHandler(service interfaces.URLShorten) *URLShortenerHandler {
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
	w.Write([]byte(host + key))
}
