package graphquery

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/Qalifah/shop/models"
	"context"
	"time"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name : "RootQuery",
	Fields : graphql.Fields{

		"getUser" : &graphql.Field{
			Type: userType,
			Description: "Get single user",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type : graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				usernameQuery, ok := params.Args["username"].(string)
				if ok {
					var user = &models.User{}
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					if err := models.UsersCollection.FindOne(ctx, bson.M{"username" : usernameQuery}).Decode(user); err == nil {
						return user, nil
					}
				}
				return models.User{}, nil
			},
		},
		"allUsers" : &graphql.Field{
			Type: graphql.NewList(userType),
			Description: "List of all users",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var users = &[]models.User{}
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				cusor, err := models.UsersCollection.Find(ctx, bson.M{})
				if err == nil {
					cusor.All(ctx, users)
					return users,  nil
				}
				return &[]models.User{}, nil
			},
		},
		"getCategory" : &graphql.Field{
			Type: categoryType,
			Description : "Get a category",
			Args : graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				nameQuery, ok := params.Args["name"].(string)
				if ok {
					var category = &models.Category{}
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					if err := models.CategoryCollection.FindOne(ctx, bson.M{"name": nameQuery}).Decode(category); err == nil {
						return category, nil
					}
				}
				return models.Category{}, nil
			},
		},
		"allCategories" : &graphql.Field{
			Type: graphql.NewList(categoryType),
			Description: "Get all categories",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var categories = &[]models.Category{}
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				cusor, err := models.CategoryCollection.Find(ctx, bson.M{})
				if err == nil {
					cusor.All(ctx, categories)
					return categories, nil
				}
				return &models.Category{}, nil
			},
		},
		"getProduct" : &graphql.Field{
			Type : productType,
			Description:  "Get a product",
			Args: graphql.FieldConfigArgument{
				"id" : &graphql.ArgumentConfig{
					Type: ObjectID,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, ok := params.Args["id"].(primitive.ObjectID)
				if ok {
					var product = &models.Product{}
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					if err := models.ProductsCollection.FindOne(ctx, bson.M{"_id" : idQuery}).Decode(product); err == nil {
						return product, nil
					}
				}
				return models.Product{}, nil
			},
		},
		"allProducts" : &graphql.Field{
			Type: graphql.NewList(productType),
			Description: "Get all products",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var products = &[]models.Product{}
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				cusor, err := models.ProductsCollection.Find(ctx, bson.M{})
				if err == nil {
					cusor.All(ctx, products)
					return products, nil					
				}
				return &[]models.Product{}, nil
			},

		},
	},
})