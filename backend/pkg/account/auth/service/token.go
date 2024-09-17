package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type refreshTokenData struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}

type idTokenCustomClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}
type refreshTokenCustomClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func validateRefreshToken(tokenString string, key []byte) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(*refreshTokenCustomClaims)
	if !ok {
		return nil, errors.New("token valid but could not parse claims")
	}
	return claims, nil
}

func validateToken(tokenString string, key []byte) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(*idTokenCustomClaims)
		if !ok {
			return nil, errors.New("token valid but could not parse claims")
		}
		fmt.Println(claims)
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*idTokenCustomClaims)
	if !ok {
		return nil, errors.New("token valid but could not parse claims")
	}
	return claims, nil
}

func generateRefreshToken(username string, key []byte) (*refreshTokenData, error) {
	unixTime := time.Now()
	tokenExp := unixTime.Add(72 * time.Hour)
	tokenId, err := uuid.NewRandom()

	if err != nil {
		log.Println("Failed to generate refresh token id")
		return nil, err
	}

	claim := refreshTokenCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenId.String(),
			IssuedAt:  unixTime.Unix(),
		},
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Printf("Failed to generate refresh token %v \n", err)
		return nil, err
	}
	return &refreshTokenData{
		SS:        tokenString,
		ID:        tokenId.String(),
		ExpiresIn: tokenExp.Sub(unixTime),
	}, nil
}

func generateIDToken(username string, key []byte) (string, error) {
	claim := idTokenCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
