package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"sync"
)

var ErrNotFound = errors.New("resource not found")

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

func (r *InMemoryRepository) Save(ctx context.Context, key string, url string, correlationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[key] = url
	return nil
}

func (r *InMemoryRepository) Find(ctx context.Context, key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.store[key]
	if !exists {
		return "", fmt.Errorf("key not found %s: %w", key, ErrNotFound)
	}
	return url, nil
}

func (r *InMemoryRepository) FindByOriginalURL(ctx context.Context, URL string) (string, error) {
	for keyStore, valueStore := range r.store {
		if valueStore == URL {
			return keyStore, nil
		}
	}

	return "", nil
}

func (r *InMemoryRepository) Close() error {
	return nil
}
