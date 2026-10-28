package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// init Slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// get configuration to run app
	flags := config.ParseFlags(os.Args[1:])

	// init storage repository
	storageRepo := initRepository(flags)

	// configure handler
	handler := configureHandler(flags, storageRepo)

	// call router
	router := router.Router(handler)

	// close openned file handler
	defer func() {
		storageRepo.Close()
		if err := storageRepo.Close(); err != nil {
			slog.Error("Ошибка при закрытии файла хранилища: ", slog.Any("error", err))
		}
	}()

	if err := run(router, flags); err != nil {
		slog.Error("failed to load configuration. Error: ", slog.Any("error", err))
		os.Exit(1)
	}
}

func configureHandler(flags config.Flags, storageRepo repository.URLRepository) *handler.Handler {
	service := service.NewURLService(storageRepo)

	cfg := &handler.FlagEnv{
		FlagRunAddr:           flags.FlagRunAddr,
		FlagRunShorternerAddr: flags.FlagRunShorternerAddr,
	}

	return handler.NewHandler(service, cfg)
}

func initRepository(configParams config.Flags) repository.URLRepository {
	var repo repository.URLRepository
	var err error

	if configParams.FlagInMemoryStorage {
		repo = repository.NewInMemoryRepository(configParams)
	} else {
		if configParams.FlagFileStoragePath != "" {
			repo, err = repository.NewFileStorgeRepository(configParams)

			if err != nil {
				slog.Error("Ошибка при инициализации репозитория хранилищац файлов: ", slog.Any("error", err))
			}
		} else {
			panic("Не указан режим работы приложения")
		}
	}

	return repo
}

func run(router chi.Router, flags config.Flags) error {
	return http.ListenAndServe(flags.FlagRunAddr, router)
}
