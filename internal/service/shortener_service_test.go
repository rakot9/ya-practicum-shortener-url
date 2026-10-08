package service

import (
	"testing"
)

// TestNewStorage проверяет корректную инициализацию хранилища
func TestNewStorage(t *testing.T) {
	s := NewStorage()
	if s == nil {
		t.Fatal("NewStorage() returned nil")
	}
	if s.store == nil {
		t.Error("store map was not initialized")
	}
}

// TestSaveAndFind проверяет базовый сценарий: сохранение и последующий поиск URL
func TestSaveAndFind(t *testing.T) {
	s := NewStorage()
	targetURL := "https://yandex.practicum.ru"

	// Тестируем сохранение
	key, err := s.Save(targetURL)
	if err != nil {
		t.Fatalf("failed to save URL: %v", err)
	}

	if len(key) < 6 || len(key) > 9 {
		t.Errorf("expected key length between 6 and 9, got %d", len(key))
	}

	// Тестируем поиск существующего ключа
	foundURL, err := s.Find(key)
	if err != nil {
		t.Fatalf("failed to find URL by key: %v", err)
	}

	if foundURL != targetURL {
		t.Errorf("expected URL %q, got %q", targetURL, foundURL)
	}
}

// TestFind_NotFound проверяет поведение Find при поиске несуществующего ключа
func TestFind_NotFound(t *testing.T) {
	s := NewStorage()

	_, err := s.Find("nonexistent")
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}

	expectedErr := "key not found in memory"
	if err.Error() != expectedErr {
		t.Errorf("expected error message %q, got %q", expectedErr, err.Error())
	}
}
