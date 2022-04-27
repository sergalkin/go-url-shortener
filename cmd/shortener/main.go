package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/handlers"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

func init() {
	address := flag.String("a", config.ServerAddress(), "SERVER_ADDRESS")
	baseURL := flag.String("b", config.BaseURL(), "BASE_URL")
	fileStoragePath := flag.String("f", config.FileStoragePath(), "FILE_STORAGE_PATH")
	databaseDSN := flag.String("d", config.DatabaseDSN(), "DATABASE_DSN")
	flag.Parse()

	config.NewConfig(
		config.WithServerAddress(*address),
		config.WithBaseURL(*baseURL),
		config.WithFileStoragePath(*fileStoragePath),
		config.WithDatabaseConnection(*databaseDSN),
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Compress(5),
		middleware.Gzip,
		middleware.Cookie,
	)

	s, err := storage.NewStorage()
	if err != nil {
		fmt.Println(err)
	}
	sequence := utils.NewSequence()

	shortenHandler := handlers.NewURLShortenerHandler(service.NewURLShortenerService(s, sequence))
	expandHandler := handlers.NewURLExpandHandler(service.NewURLExpandService(s))

	db, err := storage.NewDBConnection()
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer db.Close(ctx)
	dbHandler := handlers.NewDBHandler(db)

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortenHandler.ShortenURL)
		r.Get("/{id}", expandHandler.ExpandURL)
		r.Get("/ping", dbHandler.Ping)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shortenHandler.APIShortenURL)
		r.Get("/user/urls", expandHandler.UserURLs)
	})

	server := &http.Server{
		Addr:    config.ServerAddress(),
		Handler: r,
	}

	log.Fatalln(server.ListenAndServe())
}
