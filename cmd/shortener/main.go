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
	-t
		Sets Trusted Subnet
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/grpc"
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
	trustedSubnet := flag.String("t", config.TrustedSubnet(), "TRUSTED_SUBNET")
	grpcPort := flag.String("g", config.GRPCPort(), "GRPC_PORT")
	usingJSON := flag.String("c", config.JSONConfigPath(), "CONFIG PATH")

	flag.Parse()

	c := config.NewConfig(
		config.WithServerAddress(*address),
		config.WithBaseURL(*baseURL),
		config.WithFileStoragePath(*fileStoragePath),
		config.WithDatabaseConnection(*databaseDSN),
		config.WithEnableHTTPS(*enableHTTPS),
		config.WithTrustedSubnet(*trustedSubnet),
		config.WithGRPCPort(*grpcPort),
		config.WithJSONConfig(*usingJSON),
	)

	if c.JSONConfigPath != "" {
		c.SetJSONValues()
	}

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

	shortenService := service.NewURLShortenerService(s, seq, logger)
	shortenHandler := handlers.NewURLShortenerHandler(shortenService)

	expandService := service.NewURLExpandService(s, logger)
	expandHandler := handlers.NewURLExpandHandler(expandService)

	internalService := service.NewInternalService(s, logger)
	internalHandler := handlers.NewInternalHandler(internalService)

	db, err := storage.NewDBConnection(logger, true)
	if err != nil {
		fmt.Println(err)
	}

	ctxContext, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	ctx, cancel := context.WithTimeout(ctxContext, 30*time.Second)
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
		r.Group(func(r chi.Router) {
			r.Use(middleware.TrustedSubnet)
			r.Get("/internal/stats", internalHandler.Stats)
		})
	})

	go startGRPCServer(db, internalService, shortenService, expandService)

	if config.EnableHTTPS() {
		srv := startHTTPSServer(r, stop)
		releaseResources(ctxContext, logger, srv, db)
	} else {
		srv := startHTTPServer(r, stop)
		releaseResources(ctxContext, logger, srv, db)
	}
}

// startGRPCServer - passed to gRPC server needed services and starts it.
func startGRPCServer(db storage.DB, internal service.Internal, shorten service.URLShorten, expand service.URLExpand) {
	server := grpc.NewServer(db, internal, shorten, expand)

	listen, err := net.Listen("tcp", ":"+config.GRPCPort())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("gRPC Server started.")
	if errServe := server.Serve(listen); errServe != nil {
		log.Fatal(errServe)
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
func startHTTPSServer(r *chi.Mux, stop context.CancelFunc) *http.Server {
	pwd, errPwd := exec.Command("pwd").Output()
	if errPwd != nil {
		fmt.Println(errPwd)
		stop()
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
	go func() {
		errS := server.ListenAndServeTLS(fmt.Sprintf("%s/cert.crt", path), fmt.Sprintf("%s/cert.key", path))
		if errS != nil {
			fmt.Println(errS.Error())
			stop()
		}
	}()

	return server
}

// startHTTPServer - starts HTTP server if -s flag was not provided.
func startHTTPServer(r *chi.Mux, stop context.CancelFunc) *http.Server {
	server := &http.Server{
		Addr:    config.ServerAddress(),
		Handler: r,
	}

	fmt.Println("HTTP Server started.")
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
			stop()
		}
	}()

	return server
}

// releaseResources - realising resources, stopping db connection.
func releaseResources(ctx context.Context, l *zap.Logger, srv *http.Server, db storage.DB) {
	<-ctx.Done()
	if ctx.Err() != nil {
		fmt.Printf("Error:%v\n", ctx.Err())
	}

	l.Info("The service is shutting down...")

	if db.HasNotNilConn() {
		l.Info("Closing connection with database")

		err := db.Close(ctx)
		if err != nil {
			l.Error("Could not close connection with database")
		}

		l.Info("Connection with database closed")
	}

	if err := srv.Shutdown(ctx); err != nil {
		l.Info("app error exit", zap.Error(err))
	}

	l.Info("Done")
}
