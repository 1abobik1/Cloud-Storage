package postgresql

import (
	"context"
	"database/sql"
	"log"
)

type PostgesStorage struct {
	db *sql.DB
}

func NewPostgresStorage(storagePath string) (*PostgesStorage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, err
	}

	return &PostgesStorage{db: db}, nil
}

func (p *PostgesStorage) SaveUser(ctx context.Context, email string, password []byte) (int, error) {
	const op = "storage.postgresql.SaveUser"

	query := "INSERT INTO auth_users(email, password) VALUES ($1, $2) RETURNING id"

	var id int
	err := p.db.QueryRowContext(ctx, query, email, password).Scan(&id)
	if err != nil {
		return 0, wrapPostgresErrors(err, op)
	}

	return id, nil
}

func (p *PostgesStorage) SaveRefreshToken(ctx context.Context, refreshToken string, userID int, clientIP string) error {
	const op = "storage.postgresql.SaveRefreshToken"

	query := "INSERT INTO refresh_token(token, user_id, ip) VALUES ($1, $2, $3)"

	_, err := p.db.ExecContext(ctx, query, refreshToken, userID, clientIP)
	if err != nil {
		log.Printf("Error saving in postgresql: %s, %v \n", op, err)
		return err
	}

	return nil
}
