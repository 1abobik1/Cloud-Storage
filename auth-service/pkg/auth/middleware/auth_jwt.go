package middleware

import (
	"log"
	"net/http"
	"strings"

	validation "github.com/1abobik1/Cloud-Storage/auth-service/pkg/auth/validation"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("Error: missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		// Проверяем формат Authorization заголовка
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Error: invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		accessToken := parts[1]

		claims, err := validation.ValidateToken(accessToken, secretKey)
		if err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Сохраняем user_id в контекст для последующего использования
		if userID, ok := claims["user_id"]; ok {
			c.Set("user_id", userID)
		} else {
			log.Printf("Error: user_id missing in token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

