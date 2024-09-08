package auth

import (
	"context"
	"time"
)

type AuthenService interface {
	CreateAccessToken(ctx context.Context, username string) (string, error)
	CreateRefreshToken(ctx context.Context, username, tokenId string) (string, string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	VerifyRefreshToken(ctx context.Context, token string) (string, error)

	LogoutToken(ctx context.Context, username, token string) error
	CleanToken(ctx context.Context, userId string) error
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userId, token string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userId, prevTokenID string) error
	DeleteUserRefreshToken(ctx context.Context, userId string) error

	SetBlacklistToken(ctx context.Context, userId, token string) error
	ValidateTokenInRedis(ctx context.Context, userId, token string) (bool, error)
	GetAllTokenByUserid(ctx context.Context, userId string) (map[string]string, error)
	RemoveTokenByUserid(ctx context.Context, userId string) error
}
