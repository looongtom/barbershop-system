package user

import (
	"DoAn/entity"
	"context"
)

type UserService interface {
	Login(ctx context.Context, user entity.User) (string, error)
	Register(ctx context.Context, user entity.User) (string, error)
	GetProfile(ctx context.Context, email string) (interface{}, error)
	Logout(ctx context.Context, token string) (string, error)
	ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error)
}
