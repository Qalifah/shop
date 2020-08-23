package graphquery

import (
	"github.com/Qalifah/shop/models"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"time"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name : "RootMutation", 
	Fields: graphql.Fields{
		"updateUser" : &graphql.Field{
			Type : userType,
			Description: "Update a user",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"phoneNumber": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				username, _ := params.Args["username"].(string)
				email, emailOK := params.Args["email"].(string)
				phoneNumber, phoneNumberOK := params.Args["phoneNumber"].(string)
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				if emailOK && phoneNumberOK {
					result, err := models.UsersCollection.UpdateOne(
						ctx, 
						bson.M{"username": username}, 
						bson.D{
							{"$set", bson.M{"email": email, "phoneNumber": phoneNumber}},
							// {"$set", bson.D{{Key: "email", Value: email}, {Key: "phoneNumber", Value : phoneNumber}}},
						},
					)
					if err == nil {
						user := &models.User{}
						id, _ := result.UpsertedID.(primitive.ObjectID)
						models.UsersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(user)
						return user, nil
					}
				}
				return nil, nil
			},
		},
		"deleteUser" : &graphql.Field{
			Type: graphql.String,
			Description: "Delete a user",
			Args: graphql.FieldConfigArgument{
				"username" : &graphql.ArgumentConfig{
					Type : graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				username, ok := params.Args["username"].(string)
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				if ok {
					_, err := models.UsersCollection.DeleteOne(ctx, bson.M{"username": username})
					if err == nil {
						return "Successfully deleted!", nil
					}
				}
				return "Couldn't delete, Try again!", nil
			},
		},
		"createCategory" : &graphql.Field{
			Type: categoryType,
			Description: "Create new category",
			Args: graphql.FieldConfigArgument{
				"name" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"detail" : &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				nameQuery, nameOK := params.Args["name"].(string)
				detailQuery, detailOK := params.Args["detail"].(string)
				if nameOK && detailOK {
					var checker = &models.Category{}
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					if err := models.CategoryCollection.FindOne(ctx, bson.M{"name": nameQuery}).Decode(checker); err != nil {
						category := models.Category{Name: nameQuery, Detail: detailQuery}
						createdCategory, createErr := models.CategoryCollection.InsertOne(ctx, &category)
						if createErr == nil {
							return createdCategory, nil
						}
					}
						
				}
				return models.Category{}, nil
			},
		},
		"updateCategory" : &graphql.Field{
			Type: categoryType,
			Description: "Update a category",
			Args : graphql.FieldConfigArgument{
				"name" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"detail" : &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, idOK := params.Args["id"].(primitive.ObjectID)
				nameQuery, nameOK := params.Args["name"].(string)
				detailQuery, detailOK := params.Args["detail"].(string)
				if nameOK && detailOK && idOK {
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					checker := &models.Category{}
					if err := models.CategoryCollection.FindOne(ctx, bson.M{"nameQuery": nameQuery}).Decode(checker); err != nil {
						_, err := models.CategoryCollection.UpdateOne(
							ctx, 
							bson.M{"_id": idQuery}, 
							bson.D{
								{"$set", bson.M{"name": nameQuery, "detail": detailQuery}},
							},
						)
						if err == nil {
							category := &models.Category{}
							models.CategoryCollection.FindOne(ctx, bson.M{"_id": idQuery}).Decode(category)
							return category, nil
						}
					}
				}
				return models.Category{}, nil
			},
		},
		"deleteCategory" : &graphql.Field{
			Type: graphql.String,
			Description: "Delete a category",
			Args: graphql.FieldConfigArgument{
				"name" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name, ok := params.Args["name"].(string)
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				if ok {
					_, err := models.CategoryCollection.DeleteOne(ctx, bson.M{"name": name})
					if err == nil {
						return "Successfully deleted!", nil
					}
				}
				return "Couldn't delete, Try again!", nil
			},
		},
		"createProduct" : &graphql.Field{
			Type: productType,
			Description: "Create a product",
			Args: graphql.FieldConfigArgument{
				"categoryID" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(ObjectID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"price" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"brand" : &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description" : &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"sellerID" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(ObjectID),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				categoryIDQuery, categoryOK := params.Args["categoryID"].(primitive.ObjectID)
				nameQuery, nameOK := params.Args["name"].(string)
				priceQuery, priceOK := params.Args["price"].(int)
				brandQuery, brandOK := params.Args["brand"].(string)
				descriptionQuery, descriptionOK := params.Args["description"].(string)
				sellerIDQuery, sellerOK := params.Args["sellerID"].(primitive.ObjectID)
				if categoryOK && nameOK && priceOK && brandOK && descriptionOK && sellerOK {
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					product := models.Product{CategoryID: categoryIDQuery, Name: nameQuery, Price: priceQuery, Brand: brandQuery, Description: descriptionQuery, SellerID: sellerIDQuery}
					createdProduct, createdErr := models.ProductsCollection.InsertOne(ctx, &product)
					if createdErr == nil {
						return createdProduct, nil
					}
				}
				return nil, nil
			},
		},
		"updateProduct" : &graphql.Field{
			Type: productType,
			Description: "Update a product",
			Args: graphql.FieldConfigArgument{
				"id" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(ObjectID),
				},
				"categoryID" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(ObjectID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"price" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"brand" : &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description" : &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"sellerID" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(ObjectID),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, idOK := params.Args["id"].(primitive.ObjectID)
				categoryIDQuery, categoryOK := params.Args["categoryID"].(primitive.ObjectID)
				nameQuery, nameOK := params.Args["name"].(string)
				priceQuery, priceOK := params.Args["price"].(int)
				brandQuery, brandOK := params.Args["brand"].(string)
				descriptionQuery, descriptionOK := params.Args["description"].(string)
				sellerIDQuery, sellerOK := params.Args["sellerID"].(primitive.ObjectID)
				if idOK && categoryOK && nameOK && priceOK && brandOK && descriptionOK && sellerOK {
					checker := &models.Product{}
					ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
					if err := models.ProductsCollection.FindOne(ctx, bson.M{"_id":idQuery}).Decode(checker); err != nil {
						_, err := models.ProductsCollection.UpdateOne(
							ctx, 
							bson.M{"_id": idQuery}, 
							bson.D{
								{"$set", bson.M{"categoryID": categoryIDQuery, "name": nameQuery, "price": priceQuery, "brand": brandQuery, "description": descriptionQuery, "sellerID": sellerIDQuery}},
							},
						)
						if err == nil {
							product := &models.Product{}
							models.ProductsCollection.FindOne(ctx, bson.M{"_id": idQuery}).Decode(product)
							return product, nil
						}
					}
				}
				return models.Product{}, nil
			},
		},
		"deleteProduct" : &graphql.Field{
			Type: graphql.String,
			Description: "Delete a category",
			Args: graphql.FieldConfigArgument{
				"id" : &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(ObjectID),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(primitive.ObjectID)
				ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
				if ok {
					_, err := models.ProductsCollection.DeleteOne(ctx, bson.M{"_id": id})
					if err == nil {
						return "Successfully deleted!", nil
					}
				}
				return "Couldn't delete, Try again!", nil
			},
		},
	},
})