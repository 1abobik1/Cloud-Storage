package serviceUsers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/1abobik1/Cloud-Storage/internal/config"
	lib "github.com/1abobik1/Cloud-Storage/internal/lib/jwt"
	"github.com/1abobik1/Cloud-Storage/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UsersStorageI interface {
	SaveUser(ctx context.Context, email string, password []byte) (int, error)
	SaveRefreshToken(ctx context.Context, refreshToken string, userID int, clientIP string) error
}

type userService struct {
	userStorage UsersStorageI
	cfg         config.Config
}

func NewUserService(userStorage UsersStorageI, cfg config.Config) *userService {
	return &userService{
		userStorage: userStorage,
		cfg:         cfg,
	}
}

func (s *userService) Register(ctx context.Context, email, password, clientIP string) (accessJWT string, refreshJWT string, err error) {
	const op = "service.users.Register"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error bcrypt.GenerateFromPassword: %v, location %s \n", err, op)
		return "", "", fmt.Errorf("error bcrypt.GenerateFromPassword: %w", err)
	}

	userID, err := s.userStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Printf("Warning user already exists: %v \n", err)
			return "", "", storage.ErrUserExists
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

	if err := s.userStorage.SaveRefreshToken(ctx, refreshToken, userID, clientIP); err != nil {
		fmt.Printf("Error saving refresh token: %v", err)
		return "", "", fmt.Errorf("error saving refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
