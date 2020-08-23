package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User stores info about the users
type User struct {
	ID	primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string `bson:"username" json:"username"`
	Email string `bson:"email" json:"email"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
	Password string `bson:"password" json:"password"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

// GetBSON is responsible for providing the value(s) that will actually be saved
func (u *User) GetBSON() (interface{}, error) {
	/*if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}*/
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	type my *User
	return my(u), nil
}

