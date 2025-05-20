package handler

import (
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"

    "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/domain"
    "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/dto"
    "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/middleware"
    "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/minio"
    "github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/quota"
    "github.com/1abobik1/Cloud-Storage/file_upload_service/pkg"
    "github.com/gin-gonic/gin"
)


// CreateOne загружает один файл в MinIO
// @Summary      Загрузка одного файла
// @Description  Загружает один файл в MinIO. Необходимо передавать form-data с полем `file`.
// @Tags         Files
// @Accept       mpfd
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        file formData file true "Файл для загрузки"
// @Param        mime_type formData string false "MIME-тип (если не указан, определяется автоматически)"
// @Success      200  {object}  FileResponse             "Файл успешно загружен + данные о файле"
// @Failure      400  {object}  ErrorResponse            "Некорректный запрос"
// @Failure      403  {object}  map[string]string        "Превышена квота"
// @Failure      413  {object}  map[string]string        "Файл слишком большой"
// @Failure      500  {object}  ErrorResponse            "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/one [post]
func (h *Handler) CreateOne(c *gin.Context) {
    const op = "location internal.handler.minio_handler.minio.CreateOne"

    userID, err := pkg.GetUserID(c)
    if err != nil {
        log.Printf("Errors: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
        return
    }

    file, err := c.FormFile("file")
    if err != nil {
        log.Printf("Error: %v, %s", err, op)
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  http.StatusBadRequest,
            Error:   "No file is received",
            Details: err,
        })
        return
    }

    size := file.Size
    if err := h.quotaService.CheckQuota(c, userID, size); err != nil {
        if errors.Is(err, quota.ErrQuotaExceeded) {
            c.JSON(http.StatusForbidden, gin.H{"error": "quota exceeded"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "quota check failed"})
        return
    }

    f, err := file.Open()
    if err != nil {
        log.Printf("Error: %v, %s", err, op)
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Unable to open the file",
            Details: err,
        })
        return
    }
    defer f.Close()

    fileBytes, err := io.ReadAll(f)
    if err != nil {
        if errors.Is(err, middleware.ErrFileTooLarge) {
            c.JSON(http.StatusRequestEntityTooLarge, gin.H{
                "error": fmt.Sprintf("file is too large: limit is %d bytes", middleware.MaxFileSize),
            })
        } else {
            log.Printf("Error: %v, %s", err, op)
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Status:  http.StatusInternalServerError,
                Error:   "Unable to read the file",
                Details: err,
            })
        }
        return
    }

    fileFormat := c.PostForm("mime_type")
    if fileFormat == "" {
        fileFormat = http.DetectContentType(fileBytes)
    }

    now := time.Now().UTC()
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
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Unable to save the file",
            Details: err,
        })
        return
    }

    if err := h.quotaService.AddUsage(c, userID, size); err != nil {
        if errors.Is(err, quota.ErrNoActivePlan) {
            log.Printf("[%s] AddUsage: %v", op, err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "no active plan for user"})
            return
        }
        log.Printf("[%s] AddUsage: %v", op, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    http.StatusOK,
        "message":   "File uploaded successfully",
        "file_data": fileResp,
    })
}

// CreateMany загружает несколько файлов
// @Summary      Загрузка нескольких файлов
// @Description  Загружает несколько файлов в MinIO через form-data.
// @Tags         Files
// @Accept       mpfd
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        files formData file true "Массив файлов для загрузки" collectionFormat(multi)
// @Param        mime_type formData []string false "Массив MIME-типов (по индексу)" collectionFormat(multi)
// @Success      200  {array}   FileResponse             "Файлы успешно загружены + данные о файлах"
// @Failure      400  {object}  ErrorResponse            "Некорректная форма или отсутствуют файлы"
// @Failure      403  {object}  ErrorResponse            "Превышена квота"
// @Failure      413  {object}  map[string]string        "Один из файлов слишком большой"
// @Failure      500  {object}  ErrorResponse            "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/many [post]
func (h *Handler) CreateMany(c *gin.Context) {
    const op = "location internal.handler.minio_handler.minio.CreateMany"

    userID, err := pkg.GetUserID(c)
    if err != nil {
        log.Printf("Errors: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
        return
    }

    form, err := c.MultipartForm()
    if err != nil {
        log.Printf("Error: %v, %s", err, op)
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  http.StatusBadRequest,
            Error:   "Invalid form",
            Details: err,
        })
        return
    }

    files := form.File["files"]
    if files == nil {
        log.Printf("Error: %v, %s", err, op)
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  http.StatusBadRequest,
            Error:   "No files are received",
            Details: err,
        })
        return
    }

    var totalSize int64
    for _, fh := range files {
        totalSize += fh.Size
    }

    if err := h.quotaService.CheckQuota(c.Request.Context(), userID, totalSize); err != nil {
        if errors.Is(err, quota.ErrQuotaExceeded) {
            c.JSON(http.StatusForbidden, gin.H{"error": "quota exceeded"})
        } else {
            log.Printf("[%s] CheckQuota: %v", op, err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "quota check failed"})
        }
        return
    }

    mimeTypes := c.PostFormArray("mime_type")
    data := make(map[string]domain.FileContent)

    for i, file := range files {
        f, err := file.Open()
        if err != nil {
            log.Printf("Error: %v, %s", err, op)
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Status:  http.StatusInternalServerError,
                Error:   "Unable to open the file",
                Details: err,
            })
            return
        }

        fileBytes, err := io.ReadAll(f)
        if err != nil {
            if errors.Is(err, middleware.ErrFileTooLarge) {
                c.JSON(http.StatusRequestEntityTooLarge, gin.H{
                    "error": fmt.Sprintf("file %q is too large: limit is %d bytes", file.Filename, middleware.MaxFileSize),
                })
            } else {
                log.Printf("Error: %v, %s", err, op)
                c.JSON(http.StatusInternalServerError, ErrorResponse{
                    Status:  http.StatusInternalServerError,
                    Error:   "Unable to read the file",
                    Details: err,
                })
            }
            f.Close()
            return
        }

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

    fileRespes, err := h.minioService.CreateMany(c, data, userID)
    if err != nil {
        log.Printf("Error: %v, %s", err, op)
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Unable to save the files",
            Details: err,
        })
        return
    }

    if err := h.quotaService.AddUsage(c, userID, totalSize); err != nil {
        log.Printf("[%s] AddUsage: %v", op, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    http.StatusOK,
        "message":   "Files uploaded successfully",
        "file_data": fileRespes,
    })
}

// GetOne возвращает предварительно подписанную ссылку на скачивание одного файла
// @Summary      Получение одного файла
// @Description  Возвращает пре‐подписанную ссылку на скачивание одного файла по ID и типу.
// @Tags         Files
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        id query string true "Идентификатор объекта"
// @Param        type query string true "Категория файла (photo, unknown, video, text)"
// @Success      200  {object}  FileResponse   "Ссылка на скачивание файла"
// @Failure      400  {object}  ErrorResponse  "Некорректный запрос"
// @Failure      403  {object}  ErrorResponse  "Доступ запрещён"
// @Failure      404  {object}  ErrorResponse  "Файл не найден"
// @Failure      500  {object}  ErrorResponse  "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/one [get]
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

    fileResp, err := h.minioService.GetOne(c, objectID, userID)
    if err != nil {
        log.Printf("Error: %v,  %s", err, op)

        if errors.Is(err, minio.ErrFileNotFound) {
            c.JSON(http.StatusNotFound, ErrorResponse{
                Status:  http.StatusNotFound,
                Error:   "File not found",
                Details: fmt.Sprintf("%v, file category: %s", err.Error(), objectID.FileCategory),
            })
            return
        }

        if errors.Is(err, minio.ErrForbiddenResource) {
            c.JSON(http.StatusForbidden, ErrorResponse{
                Status:  http.StatusForbidden,
                Error:   "access to the requested resource is prohibited",
                Details: err.Error(),
            })
            return
        }

        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Enable to get the object",
            Details: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    http.StatusOK,
        "message":   "File received successfully",
        "file_data": fileResp,
    })
}

// GetMany возвращает несколько файлов (список ID) в виде ссылок
// @Summary      Получение нескольких файлов
// @Description  Возвращает пре‐подписанные ссылки на скачивание нескольких файлов.
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        objectIDs body dto.ObjectIDs true "Массив идентификаторов объектов"
// @Success      200  {array}   FileResponse    "Ссылки на скачивание файлов"
// @Failure      400  {object}  ErrorResponse   "Некорректный JSON в теле запроса"
// @Failure      403  {object}  ErrorResponse   "Доступ запрещён"
// @Failure      404  {object}  ErrorResponse   "Один из файлов не найден"
// @Failure      500  {object}  ErrorResponse   "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/many [get]
func (h *Handler) GetMany(c *gin.Context) {
    const op = "location internal.handler.minio_handler.minio.GetMany"

    userID, err := pkg.GetUserID(c)
    if err != nil {
        log.Printf("Errors: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "the user's ID was not found in the token."})
        return
    }

    var objectIDs dto.ObjectIDs
    if err := c.ShouldBindJSON(&objectIDs); err != nil {
        log.Printf("Error: %v,  %s", err, op)
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  http.StatusBadRequest,
            Error:   "Invalid request body",
            Details: err,
        })
        return
    }

    log.Printf("ObjectIDsDto: %v \n", objectIDs)

    fileResp, errs := h.minioService.GetMany(c, objectIDs.ObjectIDs, userID)
    for _, err := range errs {
        if err != nil {
            log.Printf("Error: %v,  %s", err, op)

            if errors.Is(err, minio.ErrFileNotFound) {
                c.JSON(http.StatusNotFound, ErrorResponse{
                    Status:  http.StatusNotFound,
                    Error:   "File not found",
                    Details: fmt.Sprintf("%v", err.Error()),
                })
                return
            }

            if errors.Is(err, minio.ErrForbiddenResource) {
                c.JSON(http.StatusForbidden, ErrorResponse{
                    Status:  http.StatusForbidden,
                    Error:   "access to the requested resource is prohibited",
                    Details: err.Error(),
                })
                return
            }

            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Status:  http.StatusInternalServerError,
                Error:   "Enable to get many objects",
                Details: err,
            })
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    http.StatusOK,
        "message":   "Files received successfully",
        "file_data": fileResp,
    })
}

// GetAll возвращает все файлы указанной категории
// @Summary      Получение всех файлов категории
// @Description  Возвращает пре‐подписанные ссылки на скачивание всех файлов заданной категории (photo, unknown, video, text).
// @Tags         Files
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        type query string true "Категория файлов (photo, unknown, video, text)"
// @Success      200  {array}   FileResponse    "Список ссылок на все файлы категории"
// @Failure      400  {object}  map[string]string "Некорректная категория"
// @Failure      403  {object}  ErrorResponse   "Доступ запрещён"
// @Failure      404  {object}  ErrorResponse   "Файлы не найдены"
// @Failure      500  {object}  ErrorResponse   "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/all [get]
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

    fileResp, errs := h.minioService.GetAll(c, t, userID)
    for _, err := range errs {
        if err != nil {
            log.Printf("Error: %v,  %s", err, op)

            if errors.Is(err, minio.ErrFileNotFound) {
                c.JSON(http.StatusNotFound, ErrorResponse{
                    Status:  http.StatusNotFound,
                    Error:   "File not found",
                    Details: fmt.Sprintf("%v", err.Error()),
                })
                return
            }

            if errors.Is(err, minio.ErrForbiddenResource) {
                c.JSON(http.StatusForbidden, ErrorResponse{
                    Status:  http.StatusForbidden,
                    Error:   "access to the requested resource is prohibited",
                    Details: err.Error(),
                })
                return
            }

            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Status:  http.StatusInternalServerError,
                Error:   "Enable to get many objects",
                Details: err,
            })
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    http.StatusOK,
        "message":   "All Files received successfully",
        "file_data": fileResp,
    })
}

// DeleteOne удаляет один файл
// @Summary      Удаление одного файла
// @Description  Удаляет один объект из MinIO и снижает использование квоты.
// @Tags         Files
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        id query string true "Идентификатор объекта"
// @Param        type query string true "Категория файла (photo, unknown, video, text)"
// @Success      200  {object}  map[string]string   "Файл успешно удалён"
// @Failure      400  {object}  ErrorResponse       "Некорректный запрос"
// @Failure      403  {object}  ErrorResponse       "Доступ запрещён"
// @Failure      404  {object}  ErrorResponse       "Файл не найден"
// @Failure      500  {object}  ErrorResponse       "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/one [delete]
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

    size, err := h.minioService.DeleteOne(c, objectID, userID)
    if err != nil {
        log.Printf("Error: %v,  %s", err, op)

        if errors.Is(err, minio.ErrFileNotFound) {
            c.JSON(http.StatusNotFound, ErrorResponse{
                Status:  http.StatusNotFound,
                Error:   "File not found",
                Details: fmt.Sprintf("%v", err.Error()),
            })
            return
        }

        if errors.Is(err, minio.ErrForbiddenResource) {
            c.JSON(http.StatusForbidden, ErrorResponse{
                Status:  http.StatusForbidden,
                Error:   "access to the requested resource is prohibited",
                Details: err.Error(),
            })
            return
        }

        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Cannot delete the object",
            Details: err,
        })
        return
    }

    if err := h.quotaService.RemoveUsage(c, userID, size); err != nil {
        log.Printf("RemoveUsage error: %v", err)
        c.Status(http.StatusInternalServerError)
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  http.StatusOK,
        "message": "File deleted successfully",
    })
}

// DeleteMany удаляет несколько файлов
// @Summary      Удаление нескольких файлов
// @Description  Удаляет несколько объектов из MinIO, переданных в JSON-массиве, и снижает использование квоты.
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        objectIDs body dto.ObjectIDs true "Массив идентификаторов объектов"
// @Success      200  {object}  map[string]string   "Файлы успешно удалены"
// @Failure      400  {object}  ErrorResponse       "Некорректный JSON в теле запроса"
// @Failure      403  {object}  ErrorResponse       "Доступ запрещён"
// @Failure      404  {object}  ErrorResponse       "Один из файлов не найден"
// @Failure      500  {object}  ErrorResponse       "Внутренняя ошибка сервера"
// @Security     bearerAuth
// @Router       /files/many [delete]
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
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Status:  http.StatusBadRequest,
            Error:   "Invalid request body",
            Details: err,
        })
        return
    }

    log.Printf("ObjectIDsDto: %v \n", objectIDs)

    sizes, errs := h.minioService.DeleteMany(c, objectIDs.ObjectIDs, userID)
    for _, err := range errs {
        if err != nil {
            log.Printf("Error: %v,  %s", err, op)

            if errors.Is(err, minio.ErrFileNotFound) {
                c.JSON(http.StatusNotFound, ErrorResponse{
                    Status:  http.StatusNotFound,
                    Error:   "File not found",
                    Details: fmt.Sprintf("%v", err.Error()),
                })
                return
            }

            if errors.Is(err, minio.ErrForbiddenResource) {
                c.JSON(http.StatusForbidden, ErrorResponse{
                    Status:  http.StatusForbidden,
                    Error:   "access to the requested resource is prohibited",
                    Details: err.Error(),
                })
                return
            }

            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Status:  http.StatusInternalServerError,
                Error:   "Enable to get many objects",
                Details: err,
            })
            return
        }
    }

    var totalRemoved int64
    for _, sz := range sizes {
        totalRemoved += sz
    }

    if err := h.quotaService.RemoveUsage(c.Request.Context(), userID, totalRemoved); err != nil {
        log.Printf("RemoveUsage error: %v", err)
        c.Status(http.StatusInternalServerError)
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  http.StatusOK,
        "message": "Files deleted successfully",
    })
}
