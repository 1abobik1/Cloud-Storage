package main

import (
	"context"
	"log"
	"time"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/handler"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/middleware"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/minio"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.MustLoad()

	// Инициализация соединения с Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisPort,
		DB:   0,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	// Инициализация соединения с Minio
	minioClient := minio.NewMinioClient(*cfg, redisClient)
	err := minioClient.InitMinio(cfg.MinIoPort, cfg.MinIoRootUser, cfg.MinIoRootPassword, cfg.MinIoUseSSL)
	if err != nil {
		log.Fatalf("Error connecting to Minio: %v", err)
	}

	_, s := handler.NewHandler(
		minioClient,
	)

	// Инициализация маршрутизатора Gin
	r := gin.Default()
	// cors conf
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.JWTMiddleware(cfg.JWTPublicKeyPath))

	s.RegisterRoutes(r)

	// Запуск сервера Gin
	if err := r.Run(cfg.HTTPServer); err != nil {
		panic(err)
	}
}
