package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"testing"
)

// Тест для функции run без блокировки потока реальным сервером
func Test_run(t *testing.T) {
	mockRouter := chi.NewRouter()

	// // Ошибочный тест
	t.Run("should fail with invalid address", func(t *testing.T) {
		flags := config.Flags{
			FlagRunAddr: "localhost:8081",
			FlagLog:     false,
		}

		err := run(mockRouter, flags)
		if err == nil {
			t.Error("expected an error for invalid address, got nil")
		}
	})
}
