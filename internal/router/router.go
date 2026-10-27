package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	shortenerMiddleware "github.com/rakot9/ya-practicum-shortener-url/internal/middleware"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"go.uber.org/zap"
)

func Router(configParams config.Flags) chi.Router {
	r := chi.NewRouter()
	r.Use(chiMiddleware.StripSlashes)

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	r.Use(shortenerMiddleware.ZapLogger(zapLogger))
	r.Use(shortenerMiddleware.GzipMiddleware)

	var repo repository.URLRepository

	if configParams.FlagInMemoryStorage {
		zapLogger.Info("Используем репозиторий in_memory")
		repo = repository.NewInMemoryRepository(configParams)
	} else {
		if configParams.FlagFileStoragePath != "" {
			zapLogger.Info("Используем репозиторий  file_storage")
			repo = repository.NewFileStorgeRepository(configParams)
		} else {
			zapLogger.Error("Не указан режим работы приложения")
			panic("Не указан режим работы приложения")
		}
	}

	service := service.NewURLService(repo)

	cfg := &handler.FlagEnv{
		FlagRunAddr:           configParams.FlagRunAddr,
		FlagRunShorternerAddr: configParams.FlagRunShorternerAddr,
	}

	h := handler.NewHandler(service, cfg)

	r.Post("/api/shorten", h.ApiShortURL)
	r.Get("/{id}", h.PageByID)
	r.Post("/", h.MainPage)

	return r
}
