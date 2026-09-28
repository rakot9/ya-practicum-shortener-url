package main

import (
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"net/http"
)

func main() {
	router := router.Router()

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
