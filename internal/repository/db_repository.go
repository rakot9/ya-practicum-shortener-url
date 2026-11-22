package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service/txctx"
	"time"
)

var ErrShortURLConflict = errors.New("short url already exists")

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type DBStorageRepository struct {
	pool *pgxpool.Pool
}

func NewDBStorgeRepository(pool *pgxpool.Pool) (*DBStorageRepository, error) {
	return &DBStorageRepository{
		pool: pool,
	}, nil
}

func (r *DBStorageRepository) querier(ctx context.Context) Querier {
	if tx, ok := txctx.TxFromCtx(ctx); ok {
		// идём в транзакции
		return tx
	}
	// транзакции нет — работаем с пулом напрямую
	return r.pool
}

func (r *DBStorageRepository) Save(ctx context.Context, key string, URL string, correlationID string) (string, error) {
	idV4 := uuid.New()

	var CorrelationIDDBValue *string

	if correlationID != "" {
		CorrelationIDDBValue = &correlationID
	}

	query := `INSERT INTO shortener (uuid, correlation_id, short_url, original_url, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (original_url) DO UPDATE SET original_url = EXCLUDED.original_url RETURNING short_url`

	var shortURL string

	err := r.querier(ctx).QueryRow(ctx, query, idV4, CorrelationIDDBValue, key, URL, time.Now().Format("2006-01-02 15:04:05")).Scan(&shortURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Код 23505 — нарушение уникальности (pgerrcode.UniqueViolation)
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case "shortener_short_url_key":
					return shortURL, ErrShortURLConflict
				default:
					return "", fmt.Errorf("конфликт уникальности по полю %s: %w", pgErr.ConstraintName, err)
				}
			}
		}

		return "", fmt.Errorf("не удалось создать запись в таблицы сокращения url: %w", err)
	}

	return shortURL, nil
}

func (r *DBStorageRepository) Find(ctx context.Context, key string) (string, error) {
	row := r.querier(ctx).QueryRow(ctx,
		"SELECT original_url FROM shortener WHERE short_url = $1 ORDER BY created_at DESC", key)

	// готовим переменную для чтения результата
	var URL sql.NullString
	err := row.Scan(&URL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w", ErrNotFound)
		}
		return "", fmt.Errorf("%w", ErrNotFound)
	}
	if URL.Valid {
		return URL.String, nil
	}

	return "", nil
}

func (r *DBStorageRepository) FindByOriginalURL(ctx context.Context, URL string) (string, error) {
	row := r.querier(ctx).QueryRow(ctx,
		"SELECT short_url FROM shortener WHERE original_url = $1 ORDER BY created_at DESC", URL)

	// готовим переменную для чтения результата
	var shortURL sql.NullString
	err := row.Scan(&shortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w", ErrNotFound)
		}
		return "", fmt.Errorf("%w", ErrNotFound)
	}
	if shortURL.Valid {
		return shortURL.String, nil
	}

	return "", nil
}

func (r *DBStorageRepository) Close() error {
	defer r.pool.Close()
	return nil
}
