package model

import (
	"time"
)

type Shortener struct {
	ID            int       `db:"id"`
	UserID        *int      `db:"user_id"`
	UUID          string    `db:"uuid"`
	CorrelationID *string   `db:"correlation_id"`
	ShortURL      string    `db:"short_url"`
	OriginalURL   string    `db:"original_url"`
	CreatedAt     time.Time `db:"created_at"`
}

type User struct {
	ID           int       `db:"id"`
	UUID         string    `db:"uuid"`
	Email        string    `db:"email"`
	Token        string    `db:"token"`
	RefreshToken string    `db:"refresh_token"`
	Password     string    `db:"password"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiredAt    time.Time `db:"expired_at"`
}
