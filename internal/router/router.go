package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
)

func Router(flagRunAddr string, flagRunShorternerAddr string, flagLog bool) chi.Router {
	r := chi.NewRouter()

	if flagLog {
		r.Use(middleware.Logger)
	}

	flagEnv := &handler.FlagEnv{
		FlagRunAddr:           flagRunAddr,
		FlagRunShorternerAddr: flagRunShorternerAddr,
	}

	r.Get("/{id}", handler.PageById)
	r.Post("/", handler.MainPage(flagEnv))

	return r
}
