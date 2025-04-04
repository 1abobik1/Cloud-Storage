package serviceUsers_test

import (
	"context"
	"testing"

	"github.com/1abobik1/Cloud-Storage/auth-service/config"
	"github.com/1abobik1/Cloud-Storage/auth-service/internal/domain/models"
	serviceUsers "github.com/1abobik1/Cloud-Storage/auth-service/internal/service/users"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) SaveUser(ctx context.Context, email string, password []byte) (int, error) {
	args := m.Called(ctx, email, password)
	return args.Int(0), args.Error(1)
}

func (m *mockUserStorage) UpsertRefreshToken(ctx context.Context, refreshToken string, userID int, platform string) error {
	args := m.Called(ctx, refreshToken, userID, platform)
	return args.Error(0)
}

func (m *mockUserStorage) FindUser(ctx context.Context, email string) (models.UserModel, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.UserModel), args.Error(1)
}

func (m *mockUserStorage) DeleteRefreshToken(сtx context.Context, refreshToken string) error {
	args := m.Called(сtx, refreshToken)
	return args.Error(0)
}

func TestUserService_Login(t *testing.T) {
	cfg := config.Config{
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 86400,
		PrivateKeyPath:  "testkeys/private_key.pem",
		PublicKeyPath:   "testkeys/public_key.pem",
	}

	mockStorage := new(mockUserStorage)
	service := serviceUsers.NewUserService(mockStorage, cfg)

	ctx := context.Background()
	email := "test_1@mail.ru"
	password := "test_pswd_1"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	platform := "web"

	user := models.UserModel{
		ID:          1,
		Email:       email,
		Password:    hashedPassword,
		IsActivated: false,
	}

	t.Run("успешный логин", func(t *testing.T) {
		mockStorage.On("FindUser", ctx, email).Return(user, nil)
		mockStorage.On("UpsertRefreshToken", ctx, mock.Anything, user.ID, platform).Return(nil)

		accessToken, refreshToken, err := service.Login(ctx, email, password, platform)

		assert.NoError(t, err)
		assert.NotEmpty(t, refreshToken)
		assert.NotEmpty(t, accessToken)

		mockStorage.AssertExpectations(t)
	})

	t.Run("неправильный пароль", func(t *testing.T) {
		mockStorage.On("FindUser", ctx, email).Return(user, nil)

		wrongPassword := "wrong password"

		accessToken, refreshToken, err := service.Login(ctx, email, wrongPassword, platform)

		assert.Error(t, err)
		assert.ErrorIs(t, err, serviceUsers.ErrInvalidCredentials)

		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)

		mockStorage.AssertExpectations(t)
	})

}
