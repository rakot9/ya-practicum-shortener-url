package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

// Run применяет все доступные миграции.
func Run(dsn string) error {
	// Второй аргумент ("sql") — это имя папки внутри embed.FS
	driver, err := iofs.New(fs, ".")
	if err != nil {
		return fmt.Errorf("ошибка создания iofs драйвера: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, dsn)
	if err != nil {
		return fmt.Errorf("ошибка инициализации migrate: %w", err)
	}
	defer m.Close()

	// Применяем миграции, игнорируя ошибку об отсутствии изменений
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("ошибка накатывания миграций: %w", err)
	}

	return nil
}
