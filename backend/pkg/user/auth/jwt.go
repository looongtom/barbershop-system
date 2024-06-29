package auth

import (
	"DoAn/database"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

var secretKey []byte

func init() {
	secretKey = []byte(os.Getenv("SECRET_JWT"))
}

func CreateAccessToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
func CreateRefreshToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	collection := database.ConnectMongo("blacklist_token")
	var result bson.M
	err = collection.FindOne(context.TODO(), bson.M{"token": tokenString}).Decode(&result)
	if err == nil {
		return fmt.Errorf("token has been blacklisted")
	}

	return nil
}
func VerifyRefreshToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	collection := database.ConnectMongo("blacklist_token")
	var result bson.M
	err = collection.FindOne(context.TODO(), bson.M{"refresh_token": tokenString}).Decode(&result)
	if err == nil {
		return fmt.Errorf("token has been blacklisted")
	}

	return nil

}
func RefreshToken(tokenString string) (*string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		accessToken, err := CreateAccessToken(claims.Subject)
		if err != nil {
			return nil, err
		}

		return &accessToken, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func GetSubjectFromToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	tokenString = tokenString[len("Bearer "):]

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
