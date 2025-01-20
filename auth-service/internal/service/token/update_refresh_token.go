package serviceToken

import (
	"fmt"
	"log"

	"github.com/1abobik1/Cloud-Storage/internal/utils"
)

func (s *tokenService) UpdateRefreshToken(refreshToken string, userID int) (string, error) {
	const op = "service.token.refresh.UpdateRefreshToken"

	newRefreshToken, err := utils.CreateRefreshToken(userID, s.cfg.RefreshTokenTTL, s.cfg.RefreshTokenSecretKey)
	if err != nil {
		log.Printf("Error creating access token: %v, location %s \n", err, op)
		return "", fmt.Errorf("error creating access token: %w", err)
	}

	if err := s.tokenStorage.UpdateRefreshToken(refreshToken, newRefreshToken); err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	return newRefreshToken, err
}
