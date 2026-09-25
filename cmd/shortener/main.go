package main

import (
	"net/http"
)

func main() {
	http.HandleFunc(`/{id}`, handler.pageById)
	http.HandleFunc(`/`, handler.mainPage)

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
