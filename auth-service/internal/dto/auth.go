package dto

type AuthDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Platform string `json:"platform"`
}
