package service

import (
	"DoAn/pkg/account/auth"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthServiceStruct struct {
	repo      auth.TokenRepository
	secretKey []byte
}

func NewAuthService(cacheRepo auth.TokenRepository, secretKey []byte) auth.AuthenService {
	return &AuthServiceStruct{
		repo:      cacheRepo,
		secretKey: secretKey,
	}
}

func (a AuthServiceStruct) VerifyToken(ctx context.Context, tokenString string) (string, error) {
	claims, err := validateToken(tokenString, a.secretKey)
	if err != nil {
		return "", err
	}
	isValid, err := a.repo.ValidateTokenInRedis(ctx, claims.Username, tokenString)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", errors.New("token has been blacklisted")
	}
	return claims.Username, nil
}

func (a AuthServiceStruct) LogoutToken(ctx context.Context, username, token string) error {
	err := a.repo.SetBlacklistToken(ctx, username, token)
	if err != nil {
		return err
	}
	return nil
}

func (a AuthServiceStruct) VerifyRefreshToken(ctx context.Context, token string) (string, error) {
	claims, err := validateRefreshToken(token, a.secretKey)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

func (a AuthServiceStruct) CreateRefreshToken(ctx context.Context, username, prevTokenId string) (string, string, error) {
	if prevTokenId != "" {
		err := a.repo.DeleteRefreshToken(ctx, username, prevTokenId)
		if err != nil {
			return "", "", err
		}
	}
	refreshToken, err := generateRefreshToken(username, a.secretKey)
	if err != nil {
		return "", "", err
	}

	err = a.repo.SetRefreshToken(ctx, username, refreshToken.SS, 72*time.Hour)
	if err != nil {
		return "", "", err
	}
	return refreshToken.SS, refreshToken.ID, nil
}

func (a AuthServiceStruct) CreateAccessToken(ctx context.Context, username string) (string, error) {
	claims := idTokenCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(a.secretKey)

	//accessToken, err := generateIDToken(username, a.secretKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
