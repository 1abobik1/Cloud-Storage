package serviceToken

import (
	"context"
	"fmt"
	"log"

	"github.com/1abobik1/Cloud-Storage/auth-service/internal/utils"
)

func (s *tokenService) UpdateAccessToken(refreshToken string) (string, error) {
	const op = "service.token.refresh.UpdateAccessToken"

	userID, err := s.tokenStorage.CheckRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	userKey, err := s.tokenStorage.GetUserKey(context.TODO(), userID)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	newAccessToken, err := utils.CreateAccessToken(userKey, userID, s.cfg.AccessTokenTTL, s.cfg.PrivateKeyPath)
	if err != nil {
		log.Printf("Error creating access token: %v, location %s \n", err, op)
		return "", fmt.Errorf("error creating access token: %w", err)
	}

	return newAccessToken, err
}
