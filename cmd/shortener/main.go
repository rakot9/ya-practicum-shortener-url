package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/config"
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"net/http"
)

func main() {
	config.ParseFlags()

	router := router.Router(config.FlagRunAddr, config.FlagRunShorternerAddr)

	if err := run(router); err != nil {
		panic(err)
	}
}

func run(router chi.Router) error {
	fmt.Println("Running server on", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, router)
}
