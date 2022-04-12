package storage

import (
	"sync"
)

// if Memory struct will no longer complains with Storage interface, code will be broken on building stage
var _ Storage = (*Memory)(nil)

type Memory struct {
	mu   sync.Mutex
	urls map[string]string
}

func NewMemory() *Memory {
	return &Memory{
		urls: map[string]string{},
	}
}

func (m *Memory) Store(key string, url string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[key] = url
}

func (m *Memory) Get(key string) (string, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	originalURL, ok := m.urls[key]
	return originalURL, ok
}
