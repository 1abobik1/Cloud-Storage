package main

import (
	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	dependapp "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/app"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/logger"
	miniostorage "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/storage/minIO"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting app...")

	minioClient, err := miniostorage.NewMinIOClient(cfg.MinIoPort, cfg.AccessKeyID, cfg.SecretAccessKey, "")
	if err != nil {
		panic(err)
	}

	dependapp.New(log, cfg.GRPCPort, cfg.StoragePath, cfg.PublicKeyPath, minioClient)

}
