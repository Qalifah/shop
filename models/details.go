package models

//AccessDetails contains metadata that is needed to lookup a user session in redis
type AccessDetails struct {
	AccessUUID string
	UserID	string
}