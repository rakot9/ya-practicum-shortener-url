package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

const base62Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Storage управляет хранением данных в памяти и генерацией ID
type Storage struct {
	mu    sync.RWMutex
	store map[string]string // ключ: key, значение: URL
}

// NewStorage создает новый экземпляр хранилища
func NewStorage() *Storage {
	return &Storage{
		store: make(map[string]string),
	}
}

func (s *Storage) generateKey() (string, error) {
	lengthBig, err := rand.Int(rand.Reader, big.NewInt(4)) // случайное число от 0 до 3
	if err != nil {
		return "", err
	}
	length := int(lengthBig.Int64()) + 6

	b := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(base62Alphabet)))

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		b[i] = base62Alphabet[idx.Int64()]
	}

	return string(b), nil
}

func (s *Storage) Save(URL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Генерируем ID и проверяем его на уникальность в памяти
	var shortID string
	var err error
	for {
		shortID, err = s.generateKey()
		if err != nil {
			return "", fmt.Errorf("failed to generate random id: %w", err)
		}
		// Если такой ID еще не занят, выходим из цикла
		if _, exists := s.store[shortID]; !exists {
			break
		}
	}

	// Сохраняем связь: key-> URL
	s.store[shortID] = URL
	return shortID, nil
}

func (s *Storage) Find(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	URL, exists := s.store[key]
	if !exists {
		return "", errors.New("key not found")
	}

	return URL, nil
}
