package entity

import "github.com/dgrijalva/jwt-go"

type Token struct {
	Token     string `json:"token"`
	User      string `json:"user"`
	CreatedAt string `json:"created_at"`
	ExpiredAt string `json:"expired_at"`
}
type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TokenRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
