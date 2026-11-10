package repository

import (
	"context"
	"database/sql"
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

func (r *DBStorageRepository) Save(key string, URL string) error {
	idV4 := uuid.New()

	query := `INSERT INTO shortener (uuid, short_url, original_url, created_at) VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(r.ctx, query, idV4, key, URL, time.Now().Format("2006-01-02 15:04:05"))
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
		return "", fmt.Errorf("%w", ErrNotFound)
	}
	if URL.Valid {
		return URL.String, nil
	}

	return "", nil
}

func (r *DBStorageRepository) Close() error {
	defer r.pool.Close()
	return nil
}
