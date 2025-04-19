package minioHandler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/domain"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/dto"
	erresponse "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/handler/errors_response"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/minio"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/pkg"
	"github.com/gin-gonic/gin"

	"net/http"
)

// CreateOne обработчик для создания одного объекта в хранилище MinIO из переданных данных.
func (h *Handler) CreateOne(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.CreateOne"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusBadRequest, erresponse.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "No file is received",
			Details: err,
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to open the file",
			Details: err,
		})
		return
	}
	defer f.Close()

	// Читаем содержимое файла в байтовый срез
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to read the file",
			Details: err,
		})
		return
	}

	// Пытаемся получить mime_type из формы. Если не передан - используем http.DetectContentType.
	fileFormat := c.PostForm("mime_type")
	if fileFormat == "" {
		fileFormat = http.DetectContentType(fileBytes)
	}

	now := time.Now().UTC()
	// Создаем структуру FileData для хранения данных файла
	fileData := domain.FileContent{
		Name:      file.Filename,
		Format:    fileFormat,
		CreatedAt: now,
		Data:      fileBytes,
	}

	log.Printf("USER-ID:%d FILE DATA... fileFormat:%s, fileName: %s, CreatedAt: %v", userID, fileFormat, file.Filename, now)

	fileResp, err := h.minioService.CreateOne(c, fileData, userID)
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to save the file",
			Details: err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "File uploaded successfully",
		"file_data": fileResp,
	})
}

// CreateMany обработчик для создания нескольких объектов в хранилище MinIO из переданных данных.
func (h *Handler) CreateMany(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.CreateMany"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	// Получаем multipart форму из запроса
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusBadRequest, erresponse.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "Invalid form",
			Details: err,
		})
		return
	}

	// Получаем файлы из формы
	files := form.File["files"]
	if files == nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusBadRequest, erresponse.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "No files are received",
			Details: err,
		})
		return
	}

	// Получаем массив mime_type, если передан
	mimeTypes := c.PostFormArray("mime_type")

	// Создаем map для хранения данных файлов
	data := make(map[string]domain.FileContent)

	// Проходим по каждому файлу в форме
	for i, file := range files {
		f, err := file.Open()
		if err != nil {
			log.Printf("Error: %v, %s", err, op)
			c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Unable to open the file",
				Details: err,
			})
			return
		}

		// Читаем содержимое файла в байтовый срез
		fileBytes, err := io.ReadAll(f)
		f.Close() // Закрываем файл сразу после чтения
		if err != nil {
			log.Printf("Error: %v, %s", err, op)
			c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Unable to read the file",
				Details: err,
			})
			return
		}

		// Если для файла передан MIME тип - используем его. Иначе определяем автоматически.
		var fileFormat string
		if i < len(mimeTypes) && mimeTypes[i] != "" {
			fileFormat = mimeTypes[i]
		} else {
			fileFormat = http.DetectContentType(fileBytes)
		}

		now := time.Now().UTC()
		data[file.Filename] = domain.FileContent{
			Name:      file.Filename,
			Format:    fileFormat,
			CreatedAt: now,
			Data:      fileBytes,
		}

		log.Printf("USER-ID:%d FILE DATA... fileFormat:%s, fileName: %s, CreatedAt: %v", userID, fileFormat, file.Filename, now)
	}

	// Сохраняем файлы в MinIO с помощью метода CreateMany
	fileRespes, err := h.minioService.CreateMany(c, data, userID)
	if err != nil {
		log.Printf("Error: %v, %s", err, op)
		c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to save the files",
			Details: err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Files uploaded successfully",
		"file_data": fileRespes,
	})
}


// GetOne обработчик для получения одного объекта из бакета Minio по его идентификатору.
func (h *Handler) GetOne(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.GetOne"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	objectID := dto.ObjectID{
		ObjID:        c.Query("id"),
		FileCategory: c.Query("type"),
	}

	log.Printf("objectID... ID:%s, userID:%d, FileCategory:%s", objectID.ObjID, userID, objectID.FileCategory)

	// Используем сервис MinIO для получения ссылки на объект
	fileResp, err := h.minioService.GetOne(c, objectID, userID)
	if err != nil {
		log.Printf("Error: %v,  %s", err, op)

		if errors.Is(err, minio.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, erresponse.ErrorResponse{
				Status:  http.StatusNotFound,
				Error:   "File not found",
				Details: fmt.Sprintf("%v, file category: %s", err.Error(), objectID.FileCategory),
			})
			return
		}

		if errors.Is(err, minio.ErrForbiddenResource) {
			c.JSON(http.StatusForbidden, erresponse.ErrorResponse{
				Status:  http.StatusForbidden,
				Error:   "access to the requested resource is prohibited",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Enable to get the object",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "File received successfully",
		"file_data":    fileResp,
	})
}

// GetMany обработчик для получения нескольких объектов из бакета Minio по их идентификаторам.
func (h *Handler) GetMany(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.GetMany"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	var objectIDs dto.ObjectIDs

	// Привязка JSON данных из запроса к переменной objectIDs
	if err := c.ShouldBindJSON(&objectIDs); err != nil {
		log.Printf("Error: %v,  %s", err, op)
		c.JSON(http.StatusBadRequest, erresponse.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body",
			Details: err,
		})
		return
	}

	log.Printf("ObjectIDsDto: %v \n", objectIDs)

	// Используем сервис MinIO для получения ссылок на объекты по их идентификаторам
	fileResp, errs := h.minioService.GetMany(c, objectIDs.ObjectIDs, userID)
	for _, err := range errs {
		if err != nil {
			log.Printf("Error: %v,  %s", err, op)

			if errors.Is(err, minio.ErrFileNotFound) {
				c.JSON(http.StatusNotFound, erresponse.ErrorResponse{
					Status:  http.StatusNotFound,
					Error:   "File not found",
					Details: fmt.Sprintf("%v", err.Error()),
				})
				return
			}

			if errors.Is(err, minio.ErrForbiddenResource) {
				c.JSON(http.StatusForbidden, erresponse.ErrorResponse{
					Status:  http.StatusForbidden,
					Error:   "access to the requested resource is prohibited",
					Details: err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Enable to get many objects",
				Details: err,
			})
			return
		}
	}

	// Возвращаем успешный ответ с URL-адресами полученных файлов
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Files received successfully",
		"file_data":    fileResp,
	})
}

// Метод GetAll для получения всех объектов из конкретного бакета Minio для конкретного пользователя
func (h *Handler) GetAll(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.GetAll"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	t := c.Query("type")

	if t != "photo" && t != "unknown" && t != "video" && t != "text" {
		log.Print("Error: the passed type in the query parameter. It can only be one of these types {photo, unknown, video, text}")
		c.JSON(http.StatusBadRequest, gin.H{"error": "the passed type in the query parameter. It can only be one of these types {photo, unknown, video, text}"})
		return
	}

	// Используем сервис MinIO для получения ссылок на объекты по их идентификаторам
	fileResp, errs := h.minioService.GetAll(c, t, userID)
	for _, err := range errs {
		if err != nil {
			log.Printf("Error: %v,  %s", err, op)

			if errors.Is(err, minio.ErrFileNotFound) {
				c.JSON(http.StatusNotFound, erresponse.ErrorResponse{
					Status:  http.StatusNotFound,
					Error:   "File not found",
					Details: fmt.Sprintf("%v", err.Error()),
				})
				return
			}

			if errors.Is(err, minio.ErrForbiddenResource) {
				c.JSON(http.StatusForbidden, erresponse.ErrorResponse{
					Status:  http.StatusForbidden,
					Error:   "access to the requested resource is prohibited",
					Details: err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Enable to get many objects",
				Details: err,
			})
			return
		}
	}

	// Возвращаем успешный ответ с URL-адресами полученных файлов
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Files received successfully",
		"file_data":    fileResp,
	})
}

// DeleteOne обработчик для удаления одного объекта из бакета Minio по его идентификатору.
func (h *Handler) DeleteOne(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.DeleteOne"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	objectID := dto.ObjectID{
		ObjID:        c.Query("id"),
		FileCategory: c.Query("type"),
	}

	log.Printf("objectID... ID:%s, userID:%d, FileCategory:%s", objectID.ObjID, userID, objectID.FileCategory)

	if err := h.minioService.DeleteOne(c, objectID, userID); err != nil {
		log.Printf("Error: %v,  %s", err, op)

		if errors.Is(err, minio.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, erresponse.ErrorResponse{
				Status:  http.StatusNotFound,
				Error:   "File not found",
				Details: fmt.Sprintf("%v", err.Error()),
			})
			return
		}

		if errors.Is(err, minio.ErrForbiddenResource) {
			c.JSON(http.StatusForbidden, erresponse.ErrorResponse{
				Status:  http.StatusForbidden,
				Error:   "access to the requested resource is prohibited",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Cannot delete the object",
			Details: err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "File deleted successfully",
	})
}

// DeleteMany обработчик для удаления нескольких объектов из бакета Minio по их идентификаторам.
func (h *Handler) DeleteMany(c *gin.Context) {
	const op = "location internal.handler.minio_handler.minio.DeleteMany"

	userID, err := pkg.GetUserID(c)
	if err != nil {
		log.Printf("Errors: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
		return
	}

	var objectIDs dto.ObjectIDs
	if err := c.ShouldBindJSON(&objectIDs); err != nil {
		log.Printf("Error: %v,  %s", err, op)
		c.JSON(http.StatusBadRequest, erresponse.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body",
			Details: err,
		})
		return
	}

	log.Printf("ObjectIDsDto: %v \n", objectIDs)

	errs := h.minioService.DeleteMany(c, objectIDs.ObjectIDs, userID)
	for _, err := range errs {
		if err != nil {
			log.Printf("Error: %v,  %s", err, op)

			if errors.Is(err, minio.ErrFileNotFound) {
				c.JSON(http.StatusNotFound, erresponse.ErrorResponse{
					Status:  http.StatusNotFound,
					Error:   "File not found",
					Details: fmt.Sprintf("%v", err.Error()),
				})
				return
			}

			if errors.Is(err, minio.ErrForbiddenResource) {
				c.JSON(http.StatusForbidden, erresponse.ErrorResponse{
					Status:  http.StatusForbidden,
					Error:   "access to the requested resource is prohibited",
					Details: err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, erresponse.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Enable to get many objects",
				Details: err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Files deleted successfully",
	})
}
