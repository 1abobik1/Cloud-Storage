package postgresql

import (
	"context"
	"database/sql"
	"log"

	"github.com/1abobik1/Cloud-Storage/internal/models"
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

// adding a new token, if there is already a token with such a platform, then simply update it in the database.
func (p *PostgesStorage) UpsertRefreshToken(ctx context.Context, refreshToken string, userID int, platform string) error {
	const op = "storage.postgresql.UpsertRefreshToken"

	query := `
        INSERT INTO refresh_token (token, user_id, platform)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, platform) DO UPDATE
        SET token = EXCLUDED.token;
    `
	
	_, err := p.db.ExecContext(ctx, query, refreshToken, userID, platform)
	if err != nil {
		log.Printf("Error saving in postgresql: %s, %v \n", op, err)
		return err
	}

	return nil
}

func (p *PostgesStorage) FindUser(ctx context.Context, email string) (models.UserModel, error) {
	const op = "storage.postgresql.FindUser"

	var userModel models.UserModel
	query := "SELECT id, email, password, is_activated FROM auth_users WHERE email = $1"
	err := p.db.QueryRowContext(ctx, query, email).Scan(&userModel.ID, &userModel.Email, &userModel.Password, &userModel.IsActivated)
	if err != nil {
		return models.UserModel{}, wrapPostgresErrors(err, op)
	}

	return userModel, nil
}
