package lib

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Структура данных для хранения пользовательских claims
type customClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Создание Access Token
func CreateAccessToken(userID int, duration time.Duration, accessSecret []byte) (string, error) {
	// Настраиваем claims для access токена
	claims := customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // Устанавливаем срок действия токена
			IssuedAt:  jwt.NewNumericDate(time.Now()),              // Время выпуска токена
		},
	}

	// Создаём новый токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString(accessSecret)
	if err != nil {
		return "", err // Возвращаем ошибку, если подпись не удалась
	}

	return tokenString, nil // Возвращаем готовый токен
}


// Создание Refresh Token
func CreateRefreshToken(userID int, duration time.Duration, refreshSecret []byte) (string, error) {
	// Настраиваем claims для refresh токена
	claims := customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // Срок действия токена
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // Время выпуска токена
		},
	}

	// Создаём новый токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString(refreshSecret)
	if err != nil {
		return "", err // Возвращаем ошибку, если подпись не удалась
	}

	return tokenString, nil // Возвращаем готовый токен
}
