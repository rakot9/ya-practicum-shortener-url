package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/config"
	"os"
	"testing"
)

// Тест для функции run без блокировки потока реальным сервером
func TestRunConfiguration(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	config.FlagRunAddr = "localhost:8080"
	config.FlagLog = false

	r := chi.NewRouter()

	errChan := make(chan error, 1)
	go func() {
		errChan <- run(r)
	}()
}
