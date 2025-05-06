package handlerUsers

import (
	"errors"
	"log"
	"net/http"

	"github.com/1abobik1/Cloud-Storage/auth-service/internal/dto"
	serviceUsers "github.com/1abobik1/Cloud-Storage/auth-service/internal/service/users"
	"github.com/1abobik1/Cloud-Storage/auth-service/internal/storage"
	"github.com/1abobik1/Cloud-Storage/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


func (h *userHandler) Login(c *gin.Context) {
	const op = "handler.http.users.Login"

	var authDTO dto.LogInDTO

	if err := c.BindJSON(&authDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		log.Printf("Error binding JSON: %v, location %s", err, op)
		return
	}

	validate := validator.New()
	if err := validate.Struct(authDTO); err != nil {
		log.Printf("Error: %s, location: %s", ErrValidationEmailOrPassword, op)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrValidationEmailOrPassword})
		return
	}

	if err := utils.ValidatePlatform(authDTO.Platform); err != nil {
		log.Printf("Error: %v, location: %s", err, op)
		c.JSON(http.StatusBadRequest, gin.H{"error": "platform not supported"})
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

	utils.SetRefreshTokenCookie(c, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
