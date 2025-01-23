package handlerUsers

import (
	"log"
	"net/http"

	"github.com/1abobik1/Cloud-Storage/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *userHandler) Logout(c *gin.Context) {
	const op = "handler.http.users.Logout"

	refreshToken, err := utils.GetRefreshTokenFromCookie(c)
	if err != nil {
		log.Printf("Error getting refresh token: %v, location: %s", err, op)
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.userService.RevokeRefreshToken(c, refreshToken); err != nil {
		log.Printf("Error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.Status(http.StatusOK)
}
