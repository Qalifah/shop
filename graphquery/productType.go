package graphquery

import (
	"github.com/graphql-go/graphql"
)

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name : "Product",
	Fields: graphql.Fields{
		"id" : &graphql.Field{
			Type: ObjectID,
		},
		"category" : &graphql.Field{
			Type : categoryType,
		},
		"name" : &graphql.Field{
			Type : graphql.String,
		},
		"price" : &graphql.Field{
			Type: graphql.Int,
		},
		"brand" : &graphql.Field{
			Type: graphql.String,
		},
		"description" : &graphql.Field{
			Type: graphql.String,
		},
		"seller" : &graphql.Field{
			Type: userType,
		},
	},
})