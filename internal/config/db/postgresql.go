package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DbStruct struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func Parse(dsn string, logger *zap.Logger) *DbStruct {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		logger.Error("Не удалось создать пул соединений: ", zap.Error(err))
	}

	dbStruct := &DbStruct{
		Pool: pool,
		Ctx:  ctx,
	}

	return dbStruct
}
