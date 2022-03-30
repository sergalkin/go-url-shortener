package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sergalkin/go-url-shortener.git/internal/app/handlers"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
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

	shortenHandler := handlers.NewURLShortenerHandler(service.NewURLShortenerService(memoryStorage, sequence))
	expandHandler := handlers.NewURLExpandHandler(service.NewURLExpandService(memoryStorage))

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortenHandler.ShortenURL)
		r.Get("/{id}", expandHandler.ExpandURL)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shortenHandler.ApiShortenURL)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatalln(server.ListenAndServe())
}
