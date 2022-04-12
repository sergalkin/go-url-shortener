package main

import (
	"flag"
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
	address := flag.String("a", "localhost:8080", "SERVER_ADDRESS")
	baseURL := flag.String("b", "http://localhost:8080", "BASE_URL")
	fileStoragePath := flag.String("f", "", "FILE_STORAGE_PATH")
	flag.Parse()

	config.NewConfig(
		config.WithServerAddress(*address),
		config.WithBaseURL(*baseURL),
		config.WithFileStoragePath(*fileStoragePath),
	)
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
