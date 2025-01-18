package handlerUsers

import (
	"errors"
	"log"
	"net/http"

	"github.com/1abobik1/Cloud-Storage/internal/dto"
	serviceUsers "github.com/1abobik1/Cloud-Storage/internal/service/users"
	"github.com/1abobik1/Cloud-Storage/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


func (h *userHandler) Login(c *gin.Context) {
	const op = "handler.http.users.Login"

	var authDTO dto.AuthDTO

	if err := c.BindJSON(&authDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		log.Printf("Error binding JSON: %v, location %s", err, op)
		return
	}

	validate := validator.New()
	if err := validate.Struct(authDTO); err != nil {
		log.Printf("Error: %s, loacation: %s", ErrValidationEmailOrPassword, op)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrValidationEmailOrPassword})
		return
	}

	accessToken, refreshToken, err := h.userService.Login(c, authDTO.Email, authDTO.Password, authDTO.Platform)
	if err != nil {
		if errors.Is(err, serviceUsers.ErrInvalidCredentials) {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "incorrect password or email"})
			return
		}
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		log.Printf("Error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		30*24*60*60, // 30days
		"/",
		"",
		false, // для продакшана поставить на true
		true,
	)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
