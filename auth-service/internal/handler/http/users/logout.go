package handlerUsers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// remove refresh token
func (h *userHandler) Logout(c *gin.Context) {
	const op = "handler.http.users.Logout"

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("Error refresh token not found: %v location %s", err, op)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	if len(refreshToken) == 0 {
		log.Printf("Error refresh token empty: %s", refreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token cannot be empty"})
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
