package txctx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) GetPool() *pgxpool.Pool {
	return m.pool
}

func (m *TxManager) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	// Уже внутри транзакции — не открываем вложенную, просто продолжаем в ней.
	if _, ok := TxFromCtx(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // после Commit это безопасный no-op

	if err := fn(WithTx(ctx, tx)); err != nil {
		return fmt.Errorf("tx body: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
