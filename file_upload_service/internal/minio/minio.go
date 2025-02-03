package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/dto"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/helpers"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	ErrForbiddenResource = errors.New("access to the requested resource is prohibited")
	ErrFileNotFound      = errors.New("file not found")
)

// Client интерфейс для взаимодействия с Minio
type Client interface {
	InitMinio(minioPort, minioRootUser, minioRootPassword string, minioUseSSL bool) error           // Метод для инициализации подключения к Minio
	CreateOne(ctx context.Context, file helpers.FileData, userID int) (string, error)               // Метод для создания одного объекта в бакете Minio
	CreateMany(ctx context.Context, data map[string]helpers.FileData, userID int) ([]string, error) // Метод для создания нескольких объектов в бакете Minio
	GetOne(ctx context.Context, objectID dto.ObjectID, userID int) (string, error)                  // Метод для получения одного объекта из бакета Minio
	GetMany(ctx context.Context, objectIDs []dto.ObjectID, userID int) ([]string, []error)          // Метод для получения нескольких объектов из бакета Minio
	DeleteOne(ctx context.Context, objectID dto.ObjectID, userID int) error                         // Метод для удаления одного объекта из бакета Minio
	DeleteMany(ctx context.Context, objectIDs []dto.ObjectID, userID int) []error                   // Метод для удаления нескольких объектов из бакета Minio
}

// minioClient реализация интерфейса MinioClient
type minioClient struct {
	mc  *minio.Client // Клиент Minio
	cfg config.Config
}

// NewMinioClient создает новый экземпляр Minio Client
func NewMinioClient(cfg config.Config) Client {
	return &minioClient{cfg: cfg} // Возвращает новый экземпляр minioClient с указанным именем бакета
}

func (m *minioClient) InitMinio(minioPort, minioRootUser, minioRootPassword string, minioUseSSL bool) error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(minioPort, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		return err
	}

	// Установка подключения Minio
	m.mc = client

	buckets := []string{"photo", "video", "text", "unknown"}

	for _, bucket := range buckets {
		exists, err := m.mc.BucketExists(ctx, bucket)
		if err != nil {
			return err
		}
		if !exists {
			err := m.mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateOne создает один объект в бакете Minio.
// Метод принимает структуру FileData, которая содержит имя файла и его данные.
// В случае успешной загрузки данных в бакет, метод возвращает nil, иначе возвращает ошибку.
// Все операции выполняются в контексте задачи.
func (m *minioClient) CreateOne(ctx context.Context, file helpers.FileData, userID int) (string, error) {

	objectName := generateFileName(userID)

	metadata := generateUserMetaData(userID)

	fileCategory := GetCategory(file.Data)

	options := minio.PutObjectOptions{
		ContentType:  file.Format,
		UserMetadata: metadata,
	}

	log.Printf("METADATA: %v", options.UserMetadata)

	// загрузка в объектное хранилище minio
	_, err := m.mc.PutObject(ctx, fileCategory, objectName, bytes.NewReader(file.Data), int64(len(file.Data)), options)
	if err != nil {
		return "", fmt.Errorf("error when creating an object %s: %v", file.Name, err)
	}

	// Получение URL для загруженного объекта
	url, err := m.mc.PresignedGetObject(ctx, fileCategory, objectName, m.cfg.MinIoURLLifeTime, nil)
	if err != nil {
		return "", fmt.Errorf("error when creating the URL for the object %s: %v", file.Name, err)
	}

	return url.String(), nil
}

// CreateMany создает несколько объектов в хранилище MinIO из переданных данных.
// Если происходит ошибка при создании объекта, метод возвращает ошибку,
// указывающую на неудачные объекты.
func (m *minioClient) CreateMany(ctx context.Context, data map[string]helpers.FileData, userID int) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Гарантирует освобождение контекста

	urlCh := make(chan string, len(data))
	errCh := make(chan error, len(data))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for objectID, file := range data {
		wg.Add(1)
		go func(objectID string, file helpers.FileData) {
			defer wg.Done()

			url, err := m.CreateOne(ctx, file, userID)
			if err != nil {
				mu.Lock()
				if firstErr == nil { // Запоминаем только первую ошибку
					firstErr = err
					cancel() // Отменяем все горутины
				}
				mu.Unlock()
				errCh <- err
				return
			}

			// Отправляем URL только если контекст не отменён
			select {
			case <-ctx.Done():
				return
			case urlCh <- url:
			}
		}(objectID, file)
	}

	// Ожидаем завершения всех горутин и закрываем каналы
	go func() {
		wg.Wait()
		close(urlCh)
		close(errCh)
	}()

	// Собираем результаты
	var urls []string
	for url := range urlCh {
		urls = append(urls, url)
	}

	// Если есть ошибка — возвращаем её
	if firstErr != nil {
		return nil, firstErr
	}

	return urls, nil
}

// GetOne получает один объект из бакета Minio по его идентификатору.
// Он принимает строку `objectID` в качестве параметра и возвращает срез байт данных объекта и ошибку, если такая возникает.
func (m *minioClient) GetOne(ctx context.Context, objectID dto.ObjectID, userID int) (string, error) {
	const op = "location internal.minio.GetOne"

	objInfo, err := m.mc.StatObject(ctx, objectID.FileCategory, objectID.ObjID, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("Error: %v, %s \n", err, op)
		return "", fmt.Errorf("error getting information about the object %s: %w", objectID.ObjID, ErrFileNotFound)
	}

	userIdStr, ok := objInfo.UserMetadata["User_id"]
	if !ok {
		log.Printf("Error: %v, %s \n", err, op)
		return "", fmt.Errorf("the user_id metadata was not found for the object %s: %w", objectID.ObjID, ErrFileNotFound)
	}

	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		return "", fmt.Errorf("error converting string number: %s to int", userIdStr)
	}

	if userIdInt != userID {
		return "", fmt.Errorf("you don't have access rights to other people's files: %w", ErrForbiddenResource)
	}

	url, err := m.mc.PresignedGetObject(ctx, objectID.FileCategory, objectID.ObjID, m.cfg.MinIoURLLifeTime, nil)
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		return "", OperationError{ObjectID: objectID.ObjID, Err: fmt.Errorf("error when getting the URL for the object %s: %w", objectID.ObjID, ErrFileNotFound)}
	}

	return url.String(), nil
}

func (m *minioClient) GetMany(ctx context.Context, objectIDs []dto.ObjectID, userID int) ([]string, []error) {
	urlCh := make(chan string, len(objectIDs))         // Канал для URL-адресов объектов
	errCh := make(chan OperationError, len(objectIDs)) // Канал для ошибок

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(objectID dto.ObjectID) {
			defer wg.Done()

			// Проверка отмены перед выполнением работы
			if ctx.Err() != nil {
				return
			}

			url, err := m.GetOne(ctx, objectID, userID)
			if err != nil {

				// Проверяем, не был ли контекст уже отменён
				select {
				case <-ctx.Done():
					return
				case errCh <- OperationError{
					ObjectID: objectID.ObjID,
					Err:      err,
				}:
				}

				cancel() // Отмена всех горутин
				return
			}

			// Отправка URL, если контекст не отменён
			select {
			case <-ctx.Done():
				return
			case urlCh <- url:
			}
		}(objectID)
	}

	// Закрытие каналов после завершения всех горутин.
	go func() {
		wg.Wait()
		close(urlCh)
		close(errCh)
	}()

	// Сбор URL-адресов объектов и ошибок.
	var urls []string
	var errs []error

	for url := range urlCh {
		urls = append(urls, url)
	}
	for opErr := range errCh {
		errs = append(errs, opErr.Err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return urls, nil
}

// DeleteOne удаляет один объект из бакета Minio по его идентификатору.
func (m *minioClient) DeleteOne(ctx context.Context, objectID dto.ObjectID, userID int) error {
	const op = "location internal.minio.DeleteOne"

	objInfo, err := m.mc.StatObject(ctx, objectID.FileCategory, objectID.ObjID, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("Error: %v, %s \n", err, op)
		return fmt.Errorf("error getting information about the object %s: %w", objectID.ObjID, ErrFileNotFound)
	}

	userIdStr, ok := objInfo.UserMetadata["User_id"]
	if !ok {
		log.Printf("Error: %v, %s \n", err, op)
		return fmt.Errorf("the user_id metadata was not found for the object %s: %w", objectID.ObjID, ErrFileNotFound)
	}

	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		return fmt.Errorf("error converting string number: %s to int", userIdStr)
	}

	if userIdInt != userID {
		return fmt.Errorf("you don't have access rights to other people's files: %w", ErrForbiddenResource)
	}

	err = m.mc.RemoveObject(ctx, objectID.FileCategory, objectID.ObjID, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("error: %v, %s", err, op)
		return OperationError{ObjectID: objectID.ObjID, Err: fmt.Errorf("couldn't delete selected file: %w", ErrFileNotFound)}
	}
	return nil
}

func (m *minioClient) DeleteMany(ctx context.Context, objectIDs []dto.ObjectID, userID int) []error {
	errCh := make(chan OperationError, len(objectIDs))
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(objectID dto.ObjectID) {
			defer wg.Done()

			// Проверяем, отменён ли контекст перед удалением
			if ctx.Err() != nil {
				return
			}

			err := m.DeleteOne(ctx, objectID, userID)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errCh <- OperationError{
					ObjectID: objectID.ObjID,
					Err:      err,
				}:
				}

				cancel()
			}
		}(objectID)
	}

	// Ожидание завершения горутин и закрытие канала ошибок
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Сбор ошибок
	var errs []error
	for opErr := range errCh {
		errs = append(errs, opErr.Err)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Генерируеn уникальное имя файла
func generateFileName(userID int) string {
	return fmt.Sprintf("%d/%s", userID, uuid.New().String())
}

func generateUserMetaData(userID int) map[string]string {
	return map[string]string{
		"User_id": fmt.Sprintf("%d", userID),
	}
}

func GetCategory(fileData []byte) string {
	// Определяем MIME-тип файла по его содержимому
	contentType := http.DetectContentType(fileData)

	switch {
	case strings.HasPrefix(contentType, "image/"):
		return "photo"
	case strings.HasPrefix(contentType, "video/"):
		return "video"
	case strings.HasPrefix(contentType, "text/"):
		return "text"
	default:
		return "unknown"
	}
}
