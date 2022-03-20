package shortener

import (
	"github.com/sergalkin/go-url-shortener.git/internal/app/interfaces"
	"io"
	"net/http"
	"strings"
)

const host = "http://localhost:8080/"

type URLShortenerHandler struct {
	service interfaces.URLService
}

func NewURLShortenerHandler(service interfaces.URLService) *URLShortenerHandler {
	return &URLShortenerHandler{
		service: service,
	}
}

func (h *URLShortenerHandler) URLHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.shortenURL(w, req)
		return
	case http.MethodGet:
		h.expandURL(w, req)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *URLShortenerHandler) shortenURL(w http.ResponseWriter, req *http.Request) {
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

func (h *URLShortenerHandler) expandURL(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.URL.Path, "/")

	originalLink, err := h.service.ExpandURL(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
