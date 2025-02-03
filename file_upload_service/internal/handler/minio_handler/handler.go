package minioHandler

import "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/minio"

type Handler struct {
	minioService minio.Client
}

func NewMinioHandler(
	minioService minio.Client,
) *Handler {
	return &Handler{
		minioService: minioService,
	}
}
