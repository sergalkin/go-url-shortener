package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
)

// if File struct will no longer complains with Storage interface, code will be broken on building stage
var _ Storage = (*fileStore)(nil)

type fileStore struct {
	mu   sync.Mutex
	urls map[string]string
	file *os.File
}

type urlRecord struct {
	Key string `json:"key"`
	URL string `json:"URL"`
}

func NewFile(fileStoragePath string) *fileStore {
	fs := fileStore{
		urls: map[string]string{},
	}

	f, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalln(err.Error())
	}

	fs.file = f

	if err = fs.loadFromFile(); err != nil {
		log.Fatalln(err.Error())
	}

	return &fs
}

func (m *fileStore) loadFromFile() error {
	if _, err := m.file.Seek(0, 0); err != nil {
		return err
	}

	r := &urlRecord{}
	scanner := bufio.NewScanner(m.file)

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

func (m *fileStore) Store(key, url string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[key] = url

	if err := m.saveToFile(key, url); err != nil {
		log.Fatalln(err.Error())
	}
}

func (m *fileStore) saveToFile(key, url string) error {
	e := json.NewEncoder(m.file)
	return e.Encode(urlRecord{key, url})
}

func (m *fileStore) Get(key string) (string, bool) {
	defer m.mu.Unlock()
	m.mu.Lock()

	originalURL, ok := m.urls[key]
	return originalURL, ok
}