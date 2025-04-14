package main

import (
	"time"

	"github.com/1abobik1/Cloud-Storage/auth-service/config"
	handlerToken "github.com/1abobik1/Cloud-Storage/auth-service/internal/handler/http/token"
	handlerUsers "github.com/1abobik1/Cloud-Storage/auth-service/internal/handler/http/users"
	serviceToken "github.com/1abobik1/Cloud-Storage/auth-service/internal/service/token"
	serviceUsers "github.com/1abobik1/Cloud-Storage/auth-service/internal/service/users"
	"github.com/1abobik1/Cloud-Storage/auth-service/internal/storage/postgresql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	postgresStorage, err := postgresql.NewPostgresStorageProd(cfg.StoragePath)
	if err != nil {
		panic("postgres connection error")
	}

	userService := serviceUsers.NewUserService(postgresStorage, *cfg)
	userHandler := handlerUsers.NewUserHandler(userService)

	tokenService := serviceToken.NewTokenService(postgresStorage, *cfg)
	tokenHandler := handlerToken.NewTokenHandler(tokenService)

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
	r.POST("/user/signup", userHandler.SignUp)
	r.POST("/user/login", userHandler.Login)
	r.POST("/user/logout", userHandler.Logout)

	r.POST("/token/update", tokenHandler.TokenUpdate)

	if err := r.Run(cfg.HTTPServer); err != nil {
		panic(err)
	}
}
