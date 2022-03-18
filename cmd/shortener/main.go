package main

import (
	"github.com/sergalkin/go-url-shortener.git/internal/app/shortener"
	storage2 "github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	storage := storage2.NewMemory()
	service := shortener.NewURLShortenerService(storage)
	handler := shortener.NewURLShortenerHandler(service)

	http.HandleFunc("/", handler.URLHandler)

	server := &http.Server{
		Addr: ":8080",
	}

	log.Fatalln(server.ListenAndServe())
}
