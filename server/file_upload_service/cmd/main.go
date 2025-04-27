package main

import (
	"context"
	"log"
	"time"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/config"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/handler"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/minio"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/quota"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize Redis
	rClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisPort,
		DB:   0,
	})
	if err := rClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis connection error: %v", err)
	}

	// Initialize MinIO client
	minioClient := minio.NewMinioClient(*cfg, rClient)
	if err := minioClient.InitMinio(cfg.MinIoPort, cfg.MinIoRootUser, cfg.MinIoRootPassword, cfg.MinIoUseSSL); err != nil {
		log.Fatalf("minio init error: %v", err)
	}

	// Initialize Quota service (Postgres)
	quotaSvc, err := quota.NewQuotaService(cfg.StoragePath)
	if err != nil {
		log.Fatalf("quota service init error: %v", err)
	}

	
	h := handler.NewHandler(minioClient, quotaSvc)

	// Create Gin router
	r := gin.Default()
	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware
	r.Use(middleware.JWTMiddleware(cfg.JWTPublicKeyPath))
	r.Use(middleware.MaxSizeMiddleware(middleware.MaxFileSize))
	r.Use(middleware.MaxStreamMiddleware(middleware.MaxFileSize))
	h.RegisterRoutes(r)

	// Start HTTP server
	if err := r.Run(cfg.HTTPServer); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

