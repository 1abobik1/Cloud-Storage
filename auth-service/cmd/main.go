package main

import (
	"github.com/1abobik1/Cloud-Storage/internal/config"
	handlerToken "github.com/1abobik1/Cloud-Storage/internal/handler/http/token"
	handlerUsers "github.com/1abobik1/Cloud-Storage/internal/handler/http/users"
	serviceToken "github.com/1abobik1/Cloud-Storage/internal/service/token"
	serviceUsers "github.com/1abobik1/Cloud-Storage/internal/service/users"
	"github.com/1abobik1/Cloud-Storage/internal/storage/postgresql"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	postgresStorage, err := postgresql.NewPostgresStorage(cfg.StoragePath)
	if err != nil {
		panic("postgres connection error")
	}

	userService := serviceUsers.NewUserService(postgresStorage, *cfg)
	userHandler := handlerUsers.NewUserHandler(userService)

	tokenService := serviceToken.NewTokenService(postgresStorage, *cfg)
	tokenHandler := handlerToken.NewTokenHandler(tokenService)

	r := gin.Default()

	r.POST("/user/signup", userHandler.SignUp)
	r.POST("/user/login", userHandler.Login)
	r.POST("/user/logout", userHandler.Logout)

	r.POST("/token/update", tokenHandler.TokenUpdate)

	r.Run(cfg.HTTPServer)

}
