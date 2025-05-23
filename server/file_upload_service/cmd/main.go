// @title           File Upload Service API
// @version         1.0
// @description     API для загрузки, получения и удаления файлов в облачном хранилище MinIO.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  bearerAuth
// @in                          header
// @name                        Authorization
// @description                 "Bearer {token}"
package main

import (
	"context"
	"log"
	"time"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	_ "github.com/1abobik1/Cloud-Storage/file_upload_service/docs" // docs генерируются swag и нужны для gin-swagger
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/handler"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/middleware"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/minio"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/quota"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Инициализация Redis
	rClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisPort,
		DB:   0,
	})
	if err := rClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis connection error: %v", err)
	}

	// Инициализация MinIO
	minioClient := minio.NewMinioClient(*cfg, rClient)
	if err := minioClient.InitMinio(cfg.MinIoPort, cfg.MinIoRootUser, cfg.MinIoRootPassword, cfg.MinIoUseSSL); err != nil {
		log.Fatalf("minio init error: %v", err)
	}

	// Инициализация QuotaService
	quotaSvc, err := quota.NewQuotaService(cfg.StoragePath)
	if err != nil {
		log.Fatalf("quota service init error: %v", err)
	}

	h := handler.NewHandler(minioClient, quotaSvc)

	// Создание Gin-роутера
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// JWT + ограничение размера
	r.Use(middleware.JWTMiddleware(cfg.JWTPublicKeyPath))
	r.Use(middleware.MaxSizeMiddleware(middleware.MaxFileSize))
	r.Use(middleware.MaxStreamMiddleware(middleware.MaxFileSize))

	// Регистрация маршрутов
	h.RegisterRoutes(r)

	// Подключение Swagger UI
	// По адресу http://localhost:8080/swagger/index.html будут доступны сгенерированные спецификации
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	if err := r.Run(cfg.HTTPServer); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
