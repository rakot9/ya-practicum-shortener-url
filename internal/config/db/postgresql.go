package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBStruct struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func Parse(dsn string, logger *zap.Logger) *DBStruct {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		logger.Error("Не удалось создать пул соединений: ", zap.Error(err))
	}

	dbStruct := &DBStruct{
		Pool: pool,
		Ctx:  ctx,
	}

	return dbStruct
}
