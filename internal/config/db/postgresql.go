package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service/txctx"
	"go.uber.org/zap"
)

func Parse(ctx context.Context, dsn string, logger *zap.Logger) *txctx.TxManager {

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		logger.Error("Не удалось создать пул соединений: ", zap.Error(err))
	}

	return txctx.NewTxManager(pool)
}
