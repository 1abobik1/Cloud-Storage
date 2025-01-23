package postgresql

import "database/sql"

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
