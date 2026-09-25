package router

import (
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	"net/http"
)

func Router() {
	http.HandleFunc(`/{id}`, handler.PageById)
	http.HandleFunc(`/`, handler.MainPage)
}
