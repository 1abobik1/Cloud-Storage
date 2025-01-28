package service

import (
	"context"
	"mime/multipart"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/domain/models"
)

type MinIOStorageI interface {
	Upload(ctx context.Context, userID int, fileID string, file multipart.File) (string, error)
}

type PostgresStorageI interface {
}

type fileUploaderService struct {
	minIOStorageI    MinIOStorageI
	postgresStorageI PostgresStorageI
}

func New(minIOStorageI MinIOStorageI, postgresStorageI PostgresStorageI) *fileUploaderService {
	return &fileUploaderService{
		minIOStorageI:    minIOStorageI,
		postgresStorageI: postgresStorageI,
	}
}

// UploadFile implements server.FileUploaderServiceI.
func (f *fileUploaderService) UploadFile(ctx context.Context, file multipart.File, fileInfo *models.FileModel) (int, error) {
	panic("unimplemented")
}

// DownloadFile implements server.FileUploaderServiceI.
func (f *fileUploaderService) DownloadFile(ctx context.Context, FileId int, fileModel models.FileModel) error {
	panic("unimplemented")
}

// GetFileInfo implements server.FileUploaderServiceI.
func (f *fileUploaderService) GetFileInfo(ctx context.Context, FileId int) (*models.FileModel, error) {
	panic("unimplemented")
}
