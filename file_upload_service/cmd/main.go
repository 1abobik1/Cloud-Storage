package main

import (
	"log"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/handler"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/middleware"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/minio"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	// Инициализация соединения с Minio
	minioClient := minio.NewMinioClient(*cfg)
	err := minioClient.InitMinio(cfg.MinIoPort, cfg.MinIoRootUser, cfg.MinIoRootPassword, cfg.MinIoUseSSL)
	if err != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
	}

	_, s := handler.NewHandler(
		minioClient,
	)

	// Инициализация маршрутизатора Gin
	r := gin.Default()
	r.Use(middleware.JWTMiddleware(cfg.JWTPublicKeyPath))

	s.RegisterRoutes(r)

	// Запуск сервера Gin
	if err := r.Run(cfg.HTTPServer); err != nil {
		panic(err)
	}
}
