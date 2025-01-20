package serviceToken

import "github.com/1abobik1/Cloud-Storage/internal/config"

type TokenStorageI interface {
	CheckRefreshToken(refreshToken string) (int, error)
	UpdateRefreshToken(oldRefreshToken, newRefreshToken string) error
}

type tokenService struct {
	tokenStorage TokenStorageI
	cfg          config.Config
}

func NewTokenService(tokenStorage TokenStorageI, cfg config.Config) *tokenService {
	return &tokenService{
		tokenStorage: tokenStorage,
		cfg:          cfg,
	}
}
