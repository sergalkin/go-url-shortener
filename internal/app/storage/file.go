package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

// if File struct will no longer complains with Storage interface, code will be broken on building stage
var _ Storage = (*fileStore)(nil)

type fileStore struct {
	mu       sync.Mutex
	urls     map[string]string
	userURLs map[string][]UserURLs
	filePath string
	logger   *zap.Logger
}

type urlRecord struct {
	Key string `json:"key"`
	URL string `json:"URL"`
}

func NewFile(fileStoragePath string, l *zap.Logger) *fileStore {
	fs := fileStore{
		urls:     map[string]string{},
		userURLs: map[string][]UserURLs{},
		filePath: fileStoragePath,
		logger:   l,
	}

	if err := fs.loadFromFile(); err != nil {
		l.Fatal(err.Error())
	}

	return &fs
}

func (m *fileStore) loadFromFile() error {
	f, err := os.OpenFile(m.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		m.logger.Fatal(err.Error())
	}
	defer f.Close()

	r := &urlRecord{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), r); err == nil {
			m.urls[r.Key] = r.URL
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (m *fileStore) Store(key *string, url string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[*key] = url

	var uuid string
	err := utils.Decode(middleware.GetUUID(), &uuid)
	if err != nil {
		m.logger.Error(err.Error(), zap.Error(err))
	}

	m.userURLs[uuid] = append(m.userURLs[uuid], UserURLs{ShortURL: *key, OriginalURL: url})

	if err := m.saveToFile(*key, url); err != nil {
		m.logger.Fatal(err.Error())
	}
}

func (m *fileStore) saveToFile(key, url string) error {
	f, err := os.OpenFile(m.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		m.logger.Fatal(err.Error())
	}
	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(urlRecord{key, url})
}

func (m *fileStore) Get(key string) (string, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	originalURL, ok := m.urls[key]
	return originalURL, ok
}

func (m *fileStore) LinksByUUID(uuid string) ([]UserURLs, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	links, ok := m.userURLs[uuid]
	return links, ok
}
