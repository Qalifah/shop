package models

import (
	"context"
	"time"
	"log"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/* var (
	ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	DBname = "shop"
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
)

func init() {
	checkError(err)
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	err = Client.Ping(ctx, readpref.Primary())
	checkError(err)
	fmt.Println("Successfully connected and pinged!")
}
 */

var (
	Database = ConnectDB().Database("shop")
	UsersCollection = Database.Collection("users")
	ProductsCollection = Database.Collection("products")
	CategoryCollection = Database.Collection("categories")
)

// ConnectDB initializes the database connection
 func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	checkError(err)
	err = client.Ping(ctx, readpref.Primary())
	checkError(err)
	fmt.Println("Successfully connected and pinged.")
	return client
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}