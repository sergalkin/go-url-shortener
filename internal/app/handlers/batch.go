package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/storage"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type BatchHandler struct {
	storage storage.DB
}

func NewBatchHandler(storage storage.DB) *BatchHandler {
	return &BatchHandler{
		storage: storage,
	}
}

func (h *BatchHandler) BatchInsert(w http.ResponseWriter, req *http.Request) {
	var uid string
	err := utils.Decode(middleware.GetUUID(), &uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := make([]storage.BatchRequest, 0)
	if err = json.Unmarshal(b, &requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	batchLinks, err := h.storage.BatchInsert(requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(&batchLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	w.Write(result)
}
