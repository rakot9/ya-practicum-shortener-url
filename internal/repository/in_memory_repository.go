package repository

import (
	"errors"
	"fmt"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"sync"
)

var ErrNotFound = errors.New("resource not found")

type URLRepository interface {
	Save(key, url string) error
	Find(key string) (string, error)
	Close() error
}

// InMemoryRepository реализует URLRepository для хранения данных в оперативной памяти.
type InMemoryRepository struct {
	mu           sync.RWMutex
	store        map[string]string
	configParams config.Flags
}

func NewInMemoryRepository(configParams config.Flags) *InMemoryRepository {
	return &InMemoryRepository{
		store:        make(map[string]string),
		configParams: configParams,
	}
}

func (r *InMemoryRepository) Save(key, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[key] = url
	return nil
}

func (r *InMemoryRepository) Find(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.store[key]
	if !exists {
		return "", fmt.Errorf("key not found %s: %w", key, ErrNotFound)
	}
	return url, nil
}

func (r *InMemoryRepository) Close() error {
	return nil
}
