package handlerUsers

import (
	"context"
)

const ErrValidationEmailOrPassword = "invalid email format or password must be at least 6 characters long"

type UserServiceI interface {
	Register(ctx context.Context, email, password, platform string) (accessJWT string, refreshJWT string, err error)
	Login(ctx context.Context, email, password, platform string) (accessJWT string, refreshJWT string, err error)
}

type userHandler struct {
	userService UserServiceI
}

func NewUserHandler(userService UserServiceI) *userHandler {
	return &userHandler{userService: userService}
}
