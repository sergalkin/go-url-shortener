package storage

import (
	"sync"

	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
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

func NewMemory(l *zap.Logger) *Memory {
	return &Memory{
		urls:     map[string]string{},
		userURLs: map[string][]UserURLs{},
		logger:   l,
	}
}

func (m *Memory) Store(key *string, url string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[*key] = url

	var uuid string
	err := utils.Decode(middleware.GetUUID(), &uuid)
	if err != nil {
		m.logger.Error(err.Error(), zap.Error(err))
	}

	m.userURLs[uuid] = append(m.userURLs[uuid], UserURLs{ShortURL: *key, OriginalURL: url})
}

func (m *Memory) Get(key string) (string, bool, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	originalURL, ok := m.urls[key]
	return originalURL, ok, false
}

func (m *Memory) LinksByUUID(uuid string) ([]UserURLs, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	userLinks, ok := m.userURLs[uuid]
	return userLinks, ok
}
