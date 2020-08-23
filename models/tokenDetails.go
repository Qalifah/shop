package models

// TokenDetails holds data about the access and refresh token
type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AccessUuid	string
	RefreshUuid	string
	AtExpires	int64
	RtExpires	int64
}