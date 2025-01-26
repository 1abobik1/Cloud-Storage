package validation

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var ErrTokenExpired = errors.New("token is expired")
var ErrTokenInvalid = errors.New("invalid token")

func ValidateToken(tokenString, publicKeyPath string) (jwt.MapClaims, error) {
	publicKey, err := getPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error to get public key in file: %s", publicKeyPath)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		return publicKey, nil
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

func getPublicKey(file string) (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
