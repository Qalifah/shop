package models 

// LoginCredentials are data needed to authenticate a user
type LoginCredentials struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}