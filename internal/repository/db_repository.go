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
	return &DBStorageRepository{
		ctx:  ctx,
		pool: pool,
	}, nil
}

func (r *DBStorageRepository) Save(key string, URL string, correlationId string) error {
	fmt.Print("Save db")
	idV4 := uuid.New()

	var CorrelationId *string

	if correlationId != "" {
		CorrelationId = &correlationId
	}

	query := `INSERT INTO shortener (uuid, correlation_id, short_url, original_url, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(r.ctx, query, idV4, CorrelationId, key, URL, time.Now().Format("2006-01-02 15:04:05"))
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
		fmt.Printf("error sql %w", err)
		return "", fmt.Errorf("%w", ErrNotFound)
	}
	if URL.Valid {
		return URL.String, nil
	}

	fmt.Print("Empty result")

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
