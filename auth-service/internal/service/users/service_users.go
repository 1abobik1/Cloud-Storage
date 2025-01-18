package serviceUsers

import (
	"context"

	"github.com/1abobik1/Cloud-Storage/internal/config"
	"github.com/1abobik1/Cloud-Storage/internal/models"
)

type UsersStorageI interface {
	SaveUser(ctx context.Context, email string, password []byte) (int, error)
	UpsertRefreshToken(ctx context.Context, refreshToken string, userID int, platform string) error
	FindUser(ctx context.Context, email string) (models.UserModel, error)
	DeleteRefreshToken(сtx context.Context, refreshToken string) error
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
