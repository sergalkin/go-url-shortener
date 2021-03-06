package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type decompressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w decompressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip - gzip middleware, that can compress and decompress gzip based
// middleware runs only if Content-Encoding = gzip || Accept-Encoding = gzip is present in header.
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Encoding") == "gzip" {
			var err error
			r.Body, err = gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(decompressWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
