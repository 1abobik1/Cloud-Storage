package handlerUsers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/1abobik1/Cloud-Storage/internal/dto"
	"github.com/1abobik1/Cloud-Storage/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserServiceI interface {
	Register(ctx context.Context, email, password, clientIP string) (accessJWT string, refreshJWT string, err error)
}

type userHandler struct {
	userService UserServiceI
}

func NewUserHandler(userService UserServiceI) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) SignUp(c *gin.Context) {
	var userDTO dto.UserDTO

	if err := c.BindJSON(&userDTO); err != nil {
		log.Printf("Error binding JSON: %v \n", err)
		c.Status(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format or password must be at least 6 characters long"})
		return
	}

	accessToken, refreshToken, err := h.userService.Register(c, userDTO.Email, userDTO.Password, c.ClientIP())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		log.Printf("Error Internal logic during user registration. Email: %s, Error: %v \n", userDTO.Email, err)
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
