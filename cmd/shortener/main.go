package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"go.uber.org/zap"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger       *zap.Logger
	configParams config.Flags
}

func main() {
	// get configuration to run app
	flags := config.ParseFlags(os.Args[1:])

	// init Slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// configure logger
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	app := &application{
		logger:       zapLogger,
		configParams: flags,
	}

	// init storage repository
	storageRepo := app.initRepository()

	// configure handler
	handler := app.configureHandler(storageRepo)

	// call router
	router := router.Router(handler, zapLogger)

	// close openned file handler
	defer func() {
		if err := storageRepo.Close(); err != nil {
			slog.Error("Ошибка при закрытии файла хранилища: ", slog.Any("error", err))
		}
	}()

	if err := app.run(router); err != nil {
		slog.Error("failed to load configuration. Error: ", slog.Any("error", err))
		os.Exit(1)
	}
}

func (app *application) configureHandler(storageRepo service.URLRepository) *handler.Handler {
	service := service.NewURLService(storageRepo)

	cfg := &handler.FlagEnv{
		FlagRunAddr:           app.configParams.FlagRunAddr,
		FlagRunShorternerAddr: app.configParams.FlagRunShorternerAddr,
	}

	return handler.NewHandler(service, cfg)
}

func (app *application) initRepository() service.URLRepository {
	var repo service.URLRepository
	var err error

	if app.configParams.FlagInMemoryStorage {
		repo = repository.NewInMemoryRepository(app.configParams)
	} else {
		if app.configParams.FlagFileStoragePath != "" {
			repo, err = repository.NewFileStorgeRepository(app.configParams)

			if err != nil {
				slog.Error("Ошибка при инициализации репозитория хранилищац файлов: ", slog.Any("error", err))
				os.Exit(1)
			}
		} else {
			app.logger.Error("Не указан режим работы приложения")
			os.Exit(1)
		}
	}

	return repo
}

func (app *application) run(router chi.Router) error {
	return http.ListenAndServe(app.configParams.FlagRunAddr, router)
}
