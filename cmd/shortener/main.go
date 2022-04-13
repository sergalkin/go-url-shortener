package main

import (
	"flag"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
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
	address := flag.String("a", config.ServerAddress(), "SERVER_ADDRESS")
	baseURL := flag.String("b", config.BaseURL(), "BASE_URL")
	fileStoragePath := flag.String("f", config.FileStoragePath(), "FILE_STORAGE_PATH")
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
	r.Use(
		chiMiddleware.Compress(5),
		middleware.Gzip,
	)

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
