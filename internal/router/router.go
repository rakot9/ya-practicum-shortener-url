package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	"github.com/rakot9/ya-practicum-shortener-url/internal/logger"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"go.uber.org/zap"
)

func Router(flagRunAddr string, flagRunShorternerAddr string, flagLog bool) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	r.Use(logger.ZapLogger(zapLogger))

	repo := repository.NewInMemoryRepository()

	service := service.NewURLService(repo)

	cfg := &handler.FlagEnv{
		FlagRunAddr:           flagRunAddr,
		FlagRunShorternerAddr: flagRunShorternerAddr,
	}

	h := handler.NewHandler(service, cfg)

	r.Post("/api/shorten", h.ApiShortURL)
	r.Get("/{id}", h.PageByID)
	r.Post("/", h.MainPage)

	return r
}
