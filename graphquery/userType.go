package graphquery

import (
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id" : &graphql.Field{
			Type: ObjectID,
		},
		"username" : &graphql.Field{
			Type : graphql.String,
		},
		"email" : &graphql.Field{
			Type: graphql.String,
		},
		"phoneNumber" : &graphql.Field{
			Type: graphql.String,
		},
		"password" : &graphql.Field{
			Type: graphql.String,
		},
		"createdAt" : &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})
