package txctx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var ErrTransaction = errors.New("transaction error")

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) GetPool() *pgxpool.Pool {
	return m.pool
}

func (m *TxManager) RunInTx(ctx context.Context, logger *zap.Logger, fn func(context.Context) error) error {
	// Уже внутри транзакции — не открываем вложенную, просто продолжаем в ней.
	if _, ok := TxFromCtx(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.pool.Begin(ctx)

	if err != nil {
		logger.Error("error begin transaction", zap.Error(err))
		return fmt.Errorf("begin tx: %w", ErrTransaction)
	}

	logger.Info("T.[transaction.begin]")

	defer tx.Rollback(ctx) // после Commit это безопасный no-op

	if err := fn(WithTx(ctx, tx)); err != nil {
		logger.Error("error function in transaction", zap.Error(err))
		return fmt.Errorf("error function in tx -> rollback: %w", ErrTransaction)
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("error commit transaction", zap.Error(err))
		return fmt.Errorf("error commit in tx: %w", ErrTransaction)
	}

	logger.Info("T.[transaction.commit]")

	return nil
}
