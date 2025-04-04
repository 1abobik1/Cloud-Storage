package pkg

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// GetUserID извлекает user_id из контекста
func GetUserID(ctx context.Context) (int, error) {
	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		return -1, fmt.Errorf("не удалось извлечь claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return -1, fmt.Errorf("не найден user_id в токене")
	}

	return int(userIDFloat), nil
}
