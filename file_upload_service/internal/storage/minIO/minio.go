package miniostorage

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIOClient(endpoint, accessKey, secretKey, bucket string) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{client: client, bucket: bucket}, nil
}

func getBucketName(userID int) (string, error) {
	if userID <= 0 {
		return "", fmt.Errorf("such a user does not exist")
	}
	return fmt.Sprintf("user-%d", userID), nil
}

func (m *MinIOClient) CreateUserBucket(ctx context.Context, userID int) error {
	bucketName, err := getBucketName(userID)
	if err != nil {
		return err
	}
	// Проверяем, существует ли bucket
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	// Создаем bucket, если его нет
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MinIOClient) Upload(ctx context.Context, userID int, fileID string, file multipart.File) (string, error) {
	bucketName, err := getBucketName(userID)
	if err != nil {
		return "", err
	}

	if err := m.CreateUserBucket(ctx, userID); err != nil {
		return "", err
	}

	// Чтение файла в буфер
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", err
	}

	// Загрузка файла
	_, err = m.client.PutObject(ctx, bucketName, fileID, buf, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return fileID, nil
}
