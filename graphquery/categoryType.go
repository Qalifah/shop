package graphquery

import (
	"github.com/graphql-go/graphql"
)

var categoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Category",
	Fields: graphql.Fields{
		"id" : &graphql.Field{
			Type: ObjectID,
		},
		"name" : &graphql.Field{
			Type: graphql.String,
		},
		"detail" : &graphql.Field{
			Type: graphql.String,
		},
	},
})