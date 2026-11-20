package txctx

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txKey struct{}

// WithTx кладёт транзакцию в контекст.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromCtx достаёт транзакцию. Второе значение — была ли она вообще запущена.
func TxFromCtx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
