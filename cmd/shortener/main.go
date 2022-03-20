package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sergalkin/go-url-shortener.git/internal/app/shortener"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := chi.NewRouter()

	memoryStorage := storage.NewMemory()
	sequence := utils.NewSequence()
	service := shortener.NewURLShortenerService(memoryStorage, sequence)
	handler := shortener.NewURLShortenerHandler(service)

	r.Get("/", handler.URLHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatalln(server.ListenAndServe())
}
