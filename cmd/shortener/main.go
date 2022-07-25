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
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/handlers"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/pkg/sequence"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
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

	setDefaultValuesForBuildInfo()
}

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer logger.Sync()

	rand.Seed(time.Now().UnixNano())

	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Compress(5),
		middleware.Gzip,
		middleware.Cookie,
	)

	s, err := storage.NewStorage(logger)
	if err != nil {
		fmt.Println(err)
	}
	seq := sequence.NewSequence()

	shortenHandler := handlers.NewURLShortenerHandler(service.NewURLShortenerService(s, seq, logger))
	expandHandler := handlers.NewURLExpandHandler(service.NewURLExpandService(s, logger))

	db, err := storage.NewDBConnection(logger, true)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer db.Close(ctx)
	dbHandler := handlers.NewDBHandler(db, logger)
	batchHandler := handlers.NewBatchHandler(db, logger)
	deleteHandler := handlers.NewURLDeleteHandler(service.NewURLDeleteService(db, logger))

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortenHandler.ShortenURL)
		r.Get("/{id}", expandHandler.ExpandURL)
		r.Get("/ping", dbHandler.Ping)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shortenHandler.APIShortenURL)
		r.Post("/shorten/batch", batchHandler.BatchInsert)
		r.Get("/user/urls", expandHandler.UserURLs)
		r.Delete("/user/urls", deleteHandler.Delete)
	})

	server := &http.Server{
		Addr:    config.ServerAddress(),
		Handler: r,
	}

	log.Fatalln(server.ListenAndServe())
}

func setDefaultValuesForBuildInfo() {
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
}
