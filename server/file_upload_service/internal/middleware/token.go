package middleware

import (
	"net/http"
	"strings"

	"github.com/1abobik1/Cloud-Storage/auth-service/pkg/auth/validation"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware проверяет JWT токен и добавляет claims в контекст
func JWTMiddleware(publicKeyPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		// Убираем "Bearer " и получаем сам токен
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Проверяем токен
		claims, err := validation.ValidateToken(tokenString, publicKeyPath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Добавляем claims в контекст Gin
		c.Set("claims", claims)

		// Продолжаем выполнение следующего middleware/обработчика
		c.Next()
	}
}
