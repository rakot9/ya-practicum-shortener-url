package main

import (
	"flag"
	// "github.com/go-chi/chi/v5"
	// "github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"os"
	"testing"
)

// Тест для функции run без блокировки потока реальным сервером
func TestRunConfiguration(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// flags := config.Flags{
	// 	FlagRunAddr:           "localhost:8080",
	// 	FlagRunShorternerAddr: "http://localhost:8080",
	// 	FlagLog:               false,
	// }

	// r := chi.NewRouter()

	// errChan := make(chan error, 1)
	// go func() {
	// 	errChan <- run(r)
	// }()
}
