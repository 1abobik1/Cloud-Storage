package validation

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var ErrTokenExpired = errors.New("token is expired")
var ErrTokenInvalid = errors.New("invalid token")

func ValidateToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) && (validationError.Errors&jwt.ValidationErrorExpired != 0) {
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				return claims, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
