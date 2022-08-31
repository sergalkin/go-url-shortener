package storage

import (
	"sync"

	"go.uber.org/zap"
)

// if Memory struct will no longer complains with Storage interface, code will be broken on building stage
var _ Storage = (*Memory)(nil)

type Memory struct {
	logger   *zap.Logger
	urls     map[string]string
	userURLs map[string][]UserURLs
	mu       sync.Mutex
}

// UserURLs - container for ShortURL and OriginalURL
type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewMemory - creates Memory struct
func NewMemory(l *zap.Logger) *Memory {
	return &Memory{
		urls:     map[string]string{},
		userURLs: map[string][]UserURLs{},
		logger:   l,
	}
}

// Store - storing provided URL in Memory using key.
func (m *Memory) Store(key *string, url string, uuid string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[*key] = url

	m.userURLs[uuid] = append(m.userURLs[uuid], UserURLs{ShortURL: *key, OriginalURL: url})
}

// Get - trying to get from Memory URL by its key.
func (m *Memory) Get(key string) (string, bool, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	originalURL, ok := m.urls[key]
	return originalURL, ok, false
}

// LinksByUUID - trying to get an array of UserURLs.
//  on success will return []UserURLs and true
//  on failure will return nil and false.
func (m *Memory) LinksByUUID(uuid string) ([]UserURLs, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	userLinks, ok := m.userURLs[uuid]
	return userLinks, ok
}

func (m *Memory) Stats() (int, int, error) {
	return len(m.urls), len(m.userURLs), nil
}
