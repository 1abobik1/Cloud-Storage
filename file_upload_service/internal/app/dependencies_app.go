package dependapp

import (
	"log/slog"

	grpcapp "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/app/grpc"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/service"
	miniostorage "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/storage/minIO"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/storage/postgresql"
)

type DependenciesApp struct {
	gRPCServ *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, publicKeyPath string, minioClient *miniostorage.MinIOClient) *DependenciesApp {
	postgresStorage, err := postgresql.NewPostgresStorage(storagePath)
	if err != nil {
		panic("postgres connection error")
	}

	fileUploaderService := service.New(minioClient, postgresStorage)
	grpcApp := grpcapp.New(log, fileUploaderService, grpcPort, publicKeyPath)

	return &DependenciesApp{
		gRPCServ: grpcApp,
	}
}
