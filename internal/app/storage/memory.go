package storage

import (
	"sync"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
)

// if Memory struct will no longer complains with Storage interface, code will be broken on building stage
var _ Storage = (*Memory)(nil)

type Memory struct {
	mu       sync.Mutex
	urls     map[string]string
	userURLs map[string][]UserURLs
}

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewMemory() *Memory {
	return &Memory{
		urls:     map[string]string{},
		userURLs: map[string][]UserURLs{},
	}
}

func (m *Memory) Store(key string, url string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[key] = url

	userUID := middleware.GetUUID()
	m.userURLs[userUID] = append(m.userURLs[userUID], UserURLs{ShortURL: key, OriginalURL: url})
}

func (m *Memory) Get(key string) (string, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	originalURL, ok := m.urls[key]
	return originalURL, ok
}

func (m *Memory) LinksByUUID(uuid string) ([]UserURLs, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	userLinks, ok := m.userURLs[uuid]
	return userLinks, ok
}
