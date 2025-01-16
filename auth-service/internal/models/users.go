package models

type UserModel struct {
	ID int
	Email string
	Password string
	IsActivated bool
}