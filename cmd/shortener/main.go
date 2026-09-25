package main

import (
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"net/http"
)

func main() {
	// Запускаем роутер
	router.Router()

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
