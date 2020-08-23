package models

import (
	"github.com/dgrijalva/jwt-go"
)

// AccessClaims is used to create a customised claim for the access token
type AccessClaims struct {
	Authorized bool
	AccessUUID	string
	UserID	string
	*jwt.StandardClaims
}

// RefreshClaims is used to create a customized claim for the refresh token
type RefreshClaims struct {
	RefreshUUID string
	UserID string
	*jwt.StandardClaims
}