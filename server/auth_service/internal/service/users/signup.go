package serviceUsers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/1abobik1/Cloud-Storage/auth-service/internal/storage"
	"github.com/1abobik1/Cloud-Storage/auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Register(ctx context.Context, email, password, userKey, platform string) (accessJWT string, refreshJWT string, er error) {
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

	if err := s.userStorage.SaveUserKey(ctx, userID, userKey); err != nil {
		log.Printf("Error failed to save user key: %v \n", err)
		return "", "", fmt.Errorf("error failed to save user key")
	}

	accessToken, err := utils.CreateAccessToken(userKey, userID, s.cfg.AccessTokenTTL, s.cfg.PrivateKeyPath)
	if err != nil {
		log.Printf("Error creating access token: %v \n", err)
		return "", "", fmt.Errorf("error creating access token: %w", err)
	}

	refreshToken, err := utils.CreateRefreshToken(userKey, userID, s.cfg.RefreshTokenTTL, s.cfg.PrivateKeyPath)
	if err != nil {
		log.Printf("Error creating refresh token: %v \n", err)
		return "", "", fmt.Errorf("error creating refresh token: %w", err)
	}

	if err := s.userStorage.UpsertRefreshToken(ctx, refreshToken, userID, platform); err != nil {
		log.Printf("Error upserting refresh token in db: %v", err)
		return "", "", fmt.Errorf("error upserting refresh token in db: %w", err)
	}


	if err := notifyQuotaService(s.cfg.QuotaServiceURL, userID, accessToken); err != nil {
		log.Printf("warning: failed to init free plan for user %d: %v", userID, err)
		return "","", fmt.Errorf("failed to init free plan for user %d: %v", userID, err)
	}

	return accessToken, refreshToken, nil
}

// notifyQuotaService делает POST /users/{id}/plan/init и вставляет Bearer-токен
func notifyQuotaService(baseURL string, userID int, accessToken string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/user/%d/plan/init", baseURL, userID)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	//  заголовок авторизации
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request to quota service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("quota service returned status %d", resp.StatusCode)
	}
	return nil
}
