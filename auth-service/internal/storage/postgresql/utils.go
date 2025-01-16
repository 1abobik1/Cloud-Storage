package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1abobik1/Cloud-Storage/internal/storage"
	"github.com/lib/pq"
)

func wrapPostgresErrors(err error, op string) error {
	if err == nil {
		return nil
	}

	// Проверяем, является ли ошибка PostgreSQL-ошибкой
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // Уникальное ограничение
			return fmt.Errorf("%w, location %s", storage.ErrUserExists, op)
		default:
			return fmt.Errorf("location %s: error %s: %w", op, pqErr.Code, err)
		}
	}

	// Проверяем на sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrUserNotFound
	}

	// Общая ошибка, если никакие условия не выполнились
	return fmt.Errorf("location %s: %w", op, err)
}
