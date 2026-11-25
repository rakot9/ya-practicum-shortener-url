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
	"github.com/rakot9/ya-practicum-shortener-url/internal/model"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service/txctx"
	"time"
)

const stube_user_passwrod = "scure_pass"

var ErrShortURLConflict = errors.New("short url already exists")

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
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

func (r *DBStorageRepository) Save(ctx context.Context, key string, URL string, correlationID string, userUUID string) (string, error) {
	idV4 := uuid.New()

	var CorrelationIDDBValue *string

	if correlationID != "" {
		CorrelationIDDBValue = &correlationID
	}

	userID, err := r.createUser(ctx, userUUID)

	if err != nil {
		return "", fmt.Errorf("error create user %w", err)
	}

	query := `INSERT INTO shortener (uuid, correlation_id, short_url, original_url, created_at, user_id) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (original_url) DO UPDATE SET original_url = EXCLUDED.original_url RETURNING short_url`

	var shortURL string

	err = r.querier(ctx).QueryRow(ctx, query, idV4, CorrelationIDDBValue, key, URL, time.Now().Format("2006-01-02 15:04:05"), userID).Scan(&shortURL)

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

		fmt.Printf("In Save error %v", err)

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

func (r *DBStorageRepository) FindOrCreateUser(ctx context.Context, uuid string, email string, password string) (*int, error) {

	query := `INSERT INTO users (uuid, email, password, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (uuid) DO UPDATE SET uuid = EXCLUDED.uuid RETURNING id`

	var userID int

	err := r.querier(ctx).QueryRow(ctx, query, uuid, email, password, time.Now().Format("2006-01-02 15:04:05")).Scan(&userID)

	if err != nil {
		fmt.Printf("Error insert user %v", err)
		return nil, fmt.Errorf("не удалось создать запись в таблице users: %w", err)
	}

	return &userID, nil
}

func (r *DBStorageRepository) FindUser(ctx context.Context, userUUID string) (*model.User, error) {
	row := r.querier(ctx).QueryRow(ctx, "SELECT * FROM users WHERE uuid = $1", userUUID)

	var user model.User

	if err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.Password, &user.Token, &user.RefreshToken, &user.CreatedAt, &user.ExpiredAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *DBStorageRepository) FindUsersURLs(ctx context.Context, userUUID string) ([]model.Shortener, error) {

	userID, err := r.createUser(ctx, userUUID)

	if err != nil {
		return nil, fmt.Errorf("error create user %w", err)
	}

	rows, err := r.querier(ctx).Query(ctx, "SELECT id, uuid, correlation_id, short_url, original_url, created_at, user_id FROM shortener WHERE user_id = $1 ORDER BY created_at DESC", userID)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer rows.Close()

	var userURLs []model.Shortener

	for rows.Next() {
		var userURL model.Shortener
		if err := rows.Scan(&userURL.ID, &userURL.UUID, &userURL.CorrelationID, &userURL.ShortURL, &userURL.OriginalURL, &userURL.CreatedAt, &userURL.UserID); err != nil {
			return nil, err
		}
		userURLs = append(userURLs, userURL)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userURLs, nil
}

func (r *DBStorageRepository) createUser(ctx context.Context, userUUID string) (*int, error) {
	new_user_uuid := userUUID
	new_user_email := new_user_uuid + "@test.ru"
	new_user_password := stube_user_passwrod

	// logger.Info("ZAP: Created new user",
	// 	zap.String("uuid", new_user_uuid),
	// 	zap.String("email", new_user_email),
	// )

	userID, err := r.FindOrCreateUser(ctx, new_user_uuid, new_user_email, new_user_password)

	if err != nil || userID == nil {
		// logger.Error("Ошибка вставки пользователя в БД")
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		// return
		return nil, fmt.Errorf("error create user %w", err)
	}

	return userID, nil
}
