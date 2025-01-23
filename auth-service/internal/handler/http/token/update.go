package handlerToken

import (
	"log"
	"net/http"

	"github.com/1abobik1/Cloud-Storage/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *tokenHandler) handleRefreshToken(refreshToken string) (string, error) {
	expired, claims, err := h.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	if expired {
		userID := int(claims["user_id"].(float64))
		return h.tokenService.UpdateRefreshToken(refreshToken, userID)
	}

	return refreshToken, nil
}

func (h *tokenHandler) TokenUpdate(c *gin.Context) {
	const op = "handler.http.token.RefreshToken"

	refreshToken, err := utils.GetRefreshTokenFromCookie(c)
	if err != nil {
		log.Printf("Error getting refresh token: %v, location: %s", err, op)
		c.Status(http.StatusUnauthorized)
		return
	}

	newRefreshToken, err := h.handleRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Error handling refresh token: %v, location: %s", err, op)
		c.Status(http.StatusUnauthorized)
		return
	}

	newAccessToken, err := h.tokenService.UpdateAccessToken(newRefreshToken)
	if err != nil {
		log.Printf("Error updating access token: %v, location: %s", err, op)
		c.Status(http.StatusUnauthorized)
		return
	}

	utils.SetRefreshTokenCookie(c, newRefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
