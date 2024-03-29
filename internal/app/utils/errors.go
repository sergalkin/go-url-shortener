package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrLinksConflict   = errors.New("url has been already stored") // an error that represents duplicate of URL in storage.
	ErrLinkIsDeleted   = errors.New("url has been deleted")        // an error that represents access to soft deleted URL.
	ErrGRPCWrongUserID = errors.New("wrong ID")
	ErrGRPCInternal    = errors.New("internal error occurred")
)

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
