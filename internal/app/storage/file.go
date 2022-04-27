package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
	"log"
	"os"
	"sync"
)

// if File struct will no longer complains with Storage interface, code will be broken on building stage
var _ Storage = (*fileStore)(nil)

type fileStore struct {
	mu       sync.Mutex
	urls     map[string]string
	userURLs map[string][]UserURLs
	filePath string
}

type urlRecord struct {
	Key string `json:"key"`
	URL string `json:"URL"`
}

func NewFile(fileStoragePath string) *fileStore {
	fs := fileStore{
		urls:     map[string]string{},
		userURLs: map[string][]UserURLs{},
		filePath: fileStoragePath,
	}

	if err := fs.loadFromFile(); err != nil {
		log.Fatalln(err.Error())
	}

	return &fs
}

func (m *fileStore) loadFromFile() error {
	f, err := os.OpenFile(m.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalln(err.Error())
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

func (m *fileStore) Store(key, url string) {
	defer m.mu.Unlock()
	m.mu.Lock()

	m.urls[key] = url

	var uuid string
	err := utils.Decode(middleware.GetUUID(), &uuid)
	if err != nil {
		fmt.Println(err)
	}

	m.userURLs[uuid] = append(m.userURLs[uuid], UserURLs{ShortURL: key, OriginalURL: url})

	if err := m.saveToFile(key, url); err != nil {
		log.Fatalln(err.Error())
	}
}

func (m *fileStore) saveToFile(key, url string) error {
	f, err := os.OpenFile(m.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalln(err.Error())
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
