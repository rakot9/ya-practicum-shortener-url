package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	shortenerMiddleware "github.com/rakot9/ya-practicum-shortener-url/internal/middleware"
	"go.uber.org/zap"
)

func Router(h *handler.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(chiMiddleware.StripSlashes)

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	r.Use(shortenerMiddleware.ZapLogger(zapLogger))
	r.Use(shortenerMiddleware.GzipMiddleware)

	r.Post("/api/shorten", h.APIShortURL)
	r.Get("/{id}", h.PageByID)
	r.Post("/", h.MainPage)

	return r
}
