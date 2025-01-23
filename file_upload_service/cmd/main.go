package main

import (
	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoad()

	postgresStorage, err := postgresql.NewPostgresStorage(cfg.StoragePath)
	if err != nil {
		panic("postgres connection error")
	}

	_ = postgresStorage
}
