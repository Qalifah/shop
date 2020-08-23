package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product stores info about the product
type Product struct {
	ID  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID primitive.ObjectID `bson:"categoryID" json:"categoryID"`
	Name string		`bson:"name" json:"name"`
	Price int	`bson:"price" json:"price"`
	Brand string	`bson:"brand" json:"brand"`
	Description string	`bson:"description" json:"description"`
	SellerID	primitive.ObjectID `bson:"sellerID" json:"sellerID"`
}