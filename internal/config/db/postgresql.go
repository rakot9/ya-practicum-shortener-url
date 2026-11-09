package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	//"time"
)

type DbStruct struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func Parse(dsn string, logger *zap.Logger) *DbStruct {

	/*
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Parse/validate DSN -> pgx conf
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			logger.Error("Некорректная DSN строка: ", zap.Error(err))
		}

		config.MaxConns = 10                      // Максимум 10 одновременных соединений
		config.MinConns = 2                       // Минимум 2 всегда открытых соединения
		config.MaxConnLifetime = 30 * time.Minute // Время жизни соединения

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			logger.Error("Не удалось создать пул соединений: ", zap.Error(err))
		}
	*/

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		//fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		logger.Error("Не удалось создать пул соединений: ", zap.Error(err))
	}

	dbStruct := &DbStruct{
		Pool: pool,
		Ctx:  ctx,
	}

	return dbStruct
}
