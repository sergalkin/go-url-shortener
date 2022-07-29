/*
Shortener - is a service that can transform URL to shortened URL and store it Memory/File/Database.

How to use:
	go run main.go [-flag]
The flags are:
	-a
		Sets SERVER_ADDRESS.
	-v
		Sets BASE_URL.
	-f
		Sets FILE_STORAGE_PATH.
	-d
		Sets DATABASE_DSN.
	-s
		If flag provided, starts server with HTTPS.
	-c
		Path to config file. Must be in .json format.
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/handlers"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/service"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/pkg/certificate"
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
	enableHTTPS := flag.Bool("s", config.EnableHTTPS(), "ENABLE_HTTPS")
	usingJSON := flag.String("c", config.JSONConfigPath(), "CONFIG PATH")

	flag.Parse()

	c := config.NewConfig(
		config.WithServerAddress(*address),
		config.WithBaseURL(*baseURL),
		config.WithFileStoragePath(*fileStoragePath),
		config.WithDatabaseConnection(*databaseDSN),
		config.WithEnableHTTPS(*enableHTTPS),
		config.WithJSONConfig(*usingJSON),
	)

	c.SetJSONValues()

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

	if config.EnableHTTPS() {
		log.Panic(startHTTPSServer(r))

	} else {
		log.Panic(startHTTPServer(r))
	}
}

// setDefaultValuesForBuildInfo - resigns buildValues to "N/A", if after flag parsing they still have zero values
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

// startHTTPSServer - starts HTTPS server if -s flag was provided.
func startHTTPSServer(r *chi.Mux) error {
	pwd, errPwd := exec.Command("pwd").Output()
	if errPwd != nil {
		return errPwd
	}

	var path string
	if !strings.Contains(string(pwd), "/cmd/shortener") {
		path = strings.TrimSuffix(string(pwd), "\n") + "/cmd/shortener"
	} else {
		path = "."
	}

	if _, err := os.Stat(fmt.Sprintf("%s/cert.key", path)); errors.Is(err, os.ErrNotExist) {
		certificate.Generate(path)
	}

	// конструируем менеджер TLS-сертификатов
	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("cache-dir"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist(config.ServerAddress()),
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: r,
		// для TLS-конфигурации используем менеджер сертификатов
		TLSConfig: manager.TLSConfig(),
	}

	fmt.Println("HTTPS Server started.")
	return server.ListenAndServeTLS(fmt.Sprintf("%s/cert.crt", path), fmt.Sprintf("%s/cert.key", path))
}

// startHTTPServer - starts HTTP server if -s flag was not provided.
func startHTTPServer(r *chi.Mux) error {
	server := &http.Server{
		Addr:    config.ServerAddress(),
		Handler: r,
	}

	fmt.Println("HTTP Server started.")
	return server.ListenAndServe()
}
