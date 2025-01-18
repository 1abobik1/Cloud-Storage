package handlerUsers

import (
	"errors"
	"log"
	"net/http"

	"github.com/1abobik1/Cloud-Storage/internal/dto"
	"github.com/1abobik1/Cloud-Storage/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *userHandler) SignUp(c *gin.Context) {
	const op = "handler.http.users.SignUp"

	var authDTO dto.AuthDTO

	if err := c.BindJSON(&authDTO); err != nil {
		log.Printf("Error binding JSON: %v location %s\n", err, op)
		c.Status(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(authDTO); err != nil {
		log.Printf("Error: %s, location: %s", ErrValidationEmailOrPassword, op)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrValidationEmailOrPassword})
		return
	}

	accessToken, refreshToken, err := h.userService.Register(c, authDTO.Email, authDTO.Password, authDTO.Platform)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		log.Printf("Error Internal logic during user registration. Email: %s, Error: %v \n", authDTO.Email, err)
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
