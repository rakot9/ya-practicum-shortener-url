package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

	query := `INSERT INTO users (name, email) VALUES ($1, $2)`

	_, err := r.pool.Exec(r.ctx, query, idV4, key, URL)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	return err
}

func (r *DBStorageRepository) Find(key string) (string, error) {
	return "", fmt.Errorf("key not found in storage file %s: %w", key, ErrNotFound)
}

func (r *DBStorageRepository) Close() error {
	defer r.pool.Close()
	return nil
}
