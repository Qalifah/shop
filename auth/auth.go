package auth

import (
	"os"
	"context"
	"time"
	"log"
	"encoding/json"
	"net/http"
	"github.com/Qalifah/shop/models"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	rdClient *redis.Client
	// collection = models.Client.Database(models.DBname).Collection("users")
	contx = context.Background()
)

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Unable to load .env file")
	}
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	rdClient = redis.NewClient(&redis.Options{
		Addr : dsn,
	})
	_, err := rdClient.Ping(contx).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Your redis client is ready!")
}

//Register adds new user to the database
func Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ResponseAuthError(w, err)
		return
	}
	user.Password = string(pass)
	var alreadyUsed = models.User{}
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	if err = models.UsersCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&alreadyUsed); err == nil {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Username already taken!"})
		return
	}
	createdUser, err := models.UsersCollection.InsertOne(ctx, user)
	if err != nil {
		ResponseAuthError(w, err)
		log.Println("Hey, the error is here!")
		return
	}
	json.NewEncoder(w).Encode(createdUser)
}

//ResponseAuthError is an util function that formats errors that occur during authentication
func ResponseAuthError(w http.ResponseWriter, err error) {
	er := models.ResponseError {
		Status : http.StatusUnauthorized, 
		Message : err.Error(),
	}
	json.NewEncoder(w).Encode(&er)
}

// Login verifies the user credentials
func Login(w http.ResponseWriter, r *http.Request) {
	var cred  = &models.LoginCredentials{}
	json.NewDecoder(r.Body).Decode(cred)
	tokens, err := Authenticate(w, cred.Username, cred.Password)
	if err != nil {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnprocessableEntity, Message : err.Error()})
		return
	}  
	json.NewEncoder(w).Encode(tokens)
}

//Authenticate checks if the user credential is present in the database
func Authenticate(w http.ResponseWriter, username, password string) (tokens map[string]string, err error) {
	var user = &models.User{}
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	if err = models.UsersCollection.FindOne(ctx, bson.M{"username":username}).Decode(user); err != nil {
		// json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Invalid authentication credentials"})
		return 
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Invalid authentication credentials"})
		return 
	}
	td, err := CreateToken(user.ID.Hex())
	if err != nil {
		// json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnprocessableEntity, Message : err.Error()})
		return
	}
	err = CreateAuth(user.ID.Hex(), td)
	if err != nil {
		// json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnprocessableEntity, Message : saveErr.Error()})
		return 
	}
	tokens = map[string]string {
		"access_token" : td.AccessToken,
		"refresh_token" : td.RefreshToken,
	}
	return
}

//Logout unauthenticate a user
func Logout(w http.ResponseWriter, r *http.Request) {
	tm, err := ExtractTokenMetadata(r)
	if err != nil {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Unauthorized!"})
		return
	}
	deleted, delErr := DeleteAuth(tm.AccessUUID)
	if delErr != nil || deleted == 0 {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Unauthorized!"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"Message" : "Successfully logged out!"})
}

// TokenAuthMiddleware searches and validate the token in requests
func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r)
		if err != nil {
			ResponseAuthError(w, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}