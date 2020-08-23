package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Category contains details about the categories
type Category struct {
	ID 	primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Detail string `bson:"detail,omitempty" json:"detail,omitempty"`
}

