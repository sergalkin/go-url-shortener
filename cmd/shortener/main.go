package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/handlers"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

func init() {
	config.NewConfig()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := chi.NewRouter()

	s := storage.NewStorage()
	sequence := utils.NewSequence()

	shortenHandler := handlers.NewURLShortenerHandler(service.NewURLShortenerService(s, sequence))
	expandHandler := handlers.NewURLExpandHandler(service.NewURLExpandService(s))

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortenHandler.ShortenURL)
		r.Get("/{id}", expandHandler.ExpandURL)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shortenHandler.APIShortenURL)
	})

	server := &http.Server{
		Addr:    config.ServerAddress(),
		Handler: r,
	}

	log.Fatalln(server.ListenAndServe())
}
