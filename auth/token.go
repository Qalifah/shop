package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/Qalifah/shop/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"
)

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Unable to load .env file")
	}
}

//CreateToken returns a data structure that contain access token and refresh token
func CreateToken(userID string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 3).Unix()
	td.RefreshUuid = uuid.NewV4().String()
	var err error
	accessClaim := &models.AccessClaims{
		Authorized: true,
		AccessUUID: td.AccessUuid,
		UserID: userID,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: td.AtExpires,
		},
	}
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim)
	td.AccessToken, err = aToken.SignedString([]byte(os.Getenv("ACCESS_SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	refreshClaim := &models.RefreshClaims{
		RefreshUUID: td.RefreshUuid,
		UserID: userID,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: td.RtExpires,
		},
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	td.RefreshToken, err = rToken.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

//CreateAuth adds tokens metadata to the redis
func CreateAuth(userID string, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	atErr := rdClient.Set(contx, td.AccessUuid, userID, at.Sub(now)).Err()
	if atErr != nil {
		return atErr
	}
	rtErr := rdClient.Set(contx, td.RefreshUuid, userID, rt.Sub(now)).Err()
	if rtErr != nil {
		return rtErr
	}
	return nil
}

//ExtractToken returns the token from the request header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

//VerifyToken checks if the token was created with the right claim and signing method
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	claim := &models.AccessClaims{}
	tokenString := ExtractToken(r)
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//TokenValid checks for validity of the token
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*models.AccessClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

//ExtractTokenMetadata returns metadatas about a token
func ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*models.AccessClaims)
	if ok && token.Valid {
		accessUUID := claims.AccessUUID
		if accessUUID == "" {
			return nil, err
		}
		userID := claims.UserID
		if userID == "" {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUUID: accessUUID,
			UserID: userID,
		}, nil
	}
	return nil, err
}

//FetchAuth checks if the accessUUID is present in redis
func FetchAuth(authD *models.AccessDetails) (string, error) {
	userID, err := rdClient.Get(contx, authD.AccessUUID).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}

//DeleteAuth removes a user's accessUUID from redis which implies that the user is now unauthenticated
func DeleteAuth(givenUUID string) (int64, error) {
	deleted, err := rdClient.Del(contx, givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// Refresh validates the refresh token and if valid, creates new tokens (access token and a new refresh token)
func Refresh(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	json.NewDecoder(r.Body).Decode(&mapToken)
	refreshToken := mapToken["refresh_token"]
	if refreshToken == "" {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnprocessableEntity, Message : "Invalid refresh token!"})
		return
	}
	claim := &models.RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET_KEY")), nil
	})
	if err != nil {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Refresh token has expired!"})
		return
	}
	if _, ok := token.Claims.(*models.RefreshClaims); !ok && !token.Valid {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : err.Error()})
		return
	}
	claim, ok := token.Claims.(*models.RefreshClaims)
	if ok && token.Valid {
		refreshUUID := claim.RefreshUUID
		if refreshUUID == "" {
			json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnprocessableEntity, Message : "Invalid Token!"})
			return
		}
		userID := claim.UserID
		if userID == "" {
			json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnprocessableEntity, Message : "Invalid Token!"})
			return
		}
		deleted, delErr := DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 {
			json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Unauthorized!"})
			return
		}
		ts, createErr := CreateToken(userID)
		if createErr != nil {
			json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusForbidden, Message : createErr.Error()})
			return
		}
		saveErr := CreateAuth(userID, ts)
		if saveErr != nil {
			json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusForbidden, Message : saveErr.Error()})
			return
		}
		tokens := map[string]string{
			"access_token" : ts.AccessToken,
			"refresh_token" : ts.RefreshToken,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tokens)
	} else {
		json.NewEncoder(w).Encode(models.ResponseError{Status: http.StatusUnauthorized, Message : "Refresh token has expired!"})
		return
	}
}