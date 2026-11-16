package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config/db"
	"github.com/rakot9/ya-practicum-shortener-url/internal/handler"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"github.com/rakot9/ya-practicum-shortener-url/migrations"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type application struct {
	logger       *zap.Logger
	configParams config.Flags
	dbStruct     *db.DBStruct
}

func main() {
	// get configuration to run app
	flags := config.ParseFlags(os.Args[1:])

	// configure logger
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	// configure DB connecton
	dbStruct := db.Parse(flags.FlagDBDSN, zapLogger)

	app := &application{
		logger:       zapLogger,
		configParams: flags,
		dbStruct:     dbStruct,
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
			app.logger.Error("Ошибка при закрытии файла хранилища: ", zap.Error(err))
		}
	}()

	if err := app.run(router); err != nil {
		app.logger.Error("Ошибка при загрузке конфигурации: ", zap.Error(err))
		os.Exit(1)
	}
}

func (app *application) configureHandler(storageRepo service.URLRepository) *handler.Handler {
	service := service.NewURLService(storageRepo)

	cfg := &handler.FlagEnv{
		FlagRunAddr:           app.configParams.FlagRunAddr,
		FlagRunShorternerAddr: app.configParams.FlagRunShorternerAddr,
	}

	return handler.NewHandler(service, cfg, app.dbStruct, app.logger)
}

func (app *application) initRepository() service.URLRepository {
	var repo service.URLRepository
	var err error

	// DB storage
	if app.configParams.FlagDBDSN != "" {
		// DB repo
		repo, err = repository.NewDBStorgeRepository(app.dbStruct.Ctx, app.dbStruct.Pool)

		if err != nil {
			app.logger.Error("Ошибка при инициализации репозитория DB: ", zap.Error(err))
			os.Exit(1)
		}

		if err := migrations.Run(app.configParams.FlagDBDSN); err != nil {
			app.logger.Error("Не удалось применить миграции", zap.Error(err))
		}

		app.logger.Info("Миграции актуальны, запускаем приложение...")
	} else {
		if app.configParams.FlagFileStoragePath != "" {
			repo, err = repository.NewFileStorgeRepository(app.configParams)

			if err != nil {
				app.logger.Error("Ошибка при инициализации репозитория хранилищац файлов: ", zap.Error(err))
				os.Exit(1)
			}
		} else {
			repo = repository.NewInMemoryRepository(app.configParams)
		}
	}

	return repo
}

func (app *application) run(router chi.Router) error {
	return http.ListenAndServe(app.configParams.FlagRunAddr, router)
}
