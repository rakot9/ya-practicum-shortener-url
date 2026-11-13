package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type DBStorageRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewDBStorgeRepository(ctx context.Context, pool *pgxpool.Pool) (*DBStorageRepository, error) {

	// Это дичь но миграции не запускаются в тесте, поэтому я включил запросы вот сюда!!!
	// SQL-запрос на создание таблицы в схеме public
	query := `
		CREATE TABLE IF NOT EXISTS shortener (
			id SERIAL PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL,
			correlation_id VARCHAR(36) NULL,
			short_url VARCHAR(255) NOT NULL,
			original_url text NOT NULL,
			created_at timestamp NOT NULL
		);
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table urls: %w", err)
	}

	createIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_original_url ON shortener(original_url);
	`
	if _, err := pool.Exec(ctx, createIndexQuery); err != nil {
		return nil, fmt.Errorf("failed to create index idx_original_url: %w", err)
	}

	return &DBStorageRepository{
		ctx:  ctx,
		pool: pool,
	}, nil
}

func (r *DBStorageRepository) Save(key string, URL string, correlationID string) error {
	idV4 := uuid.New()

	var CorrelationIDDBValue *string

	if correlationID != "" {
		CorrelationIDDBValue = &correlationID
	}

	query := `INSERT INTO shortener (uuid, correlation_id, short_url, original_url, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(r.ctx, query, idV4, CorrelationIDDBValue, key, URL, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("не удалось создать запись в таблицы сокращения url: %w", err)
	}

	return err
}

func (r *DBStorageRepository) Find(key string) (string, error) {
	row := r.pool.QueryRow(r.ctx,
		"SELECT original_url FROM shortener WHERE short_url = $1 ORDER BY created_at DESC", key)

	// готовим переменную для чтения результата
	var URL sql.NullString
	err := row.Scan(&URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w", ErrNotFound)
		}
		return "", fmt.Errorf("%w", ErrNotFound)
	}
	if URL.Valid {
		return URL.String, nil
	}

	return "", nil
}

func (r *DBStorageRepository) FindByOriginalURL(URL string) (string, error) {
	row := r.pool.QueryRow(r.ctx,
		"SELECT short_url FROM shortener WHERE original_url = $1 ORDER BY created_at DESC", URL)

	// готовим переменную для чтения результата
	var shortURL sql.NullString
	err := row.Scan(&shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
