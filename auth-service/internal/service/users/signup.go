package serviceUsers

import (
	"context"
	"errors"
	"fmt"
	"log"

	lib "github.com/1abobik1/Cloud-Storage/internal/lib/jwt"
	"github.com/1abobik1/Cloud-Storage/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Register(ctx context.Context, email, password, platform string) (accessJWT string, refreshJWT string, err error) {
	const op = "service.users.Register"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error bcrypt.GenerateFromPassword: %v, location %s \n", err, op)
		return "", "", fmt.Errorf("error bcrypt.GenerateFromPassword: %w", err)
	}

	userID, err := s.userStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Printf("Warning: %v \n", err)
			return "", "", err
		}

		log.Printf("Error failed to save user: %v \n", err)
		return "", "", err
	}

	accessToken, err := lib.CreateAccessToken(userID, s.cfg.AccessTokenTTL, s.cfg.AccessTokenSecretKey)
	if err != nil {
		log.Printf("Error creating access token: %v \n", err)
		return "", "", fmt.Errorf("error creating access token: %w", err)
	}

	refreshToken, err := lib.CreateRefreshToken(userID, s.cfg.RefreshTokenTTL, s.cfg.RefreshTokenSecretKey)
	if err != nil {
		log.Printf("Error creating refresh token: %v \n", err)
		return "", "", fmt.Errorf("error creating refresh token: %w", err)
	}

	if err := s.userStorage.UpsertRefreshToken(ctx, refreshToken, userID, platform); err != nil {
		log.Printf("Error upserting refresh token in db: %v", err)
		return "", "", fmt.Errorf("error upserting refresh token in db: %w", err)
	}
	

	return accessToken, refreshToken, nil
}
