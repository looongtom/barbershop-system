package user

import (
	"DoAn/entity"
	"context"
)

type UserRepository interface {
	Login(ctx context.Context, user entity.User) (string, error)
	Register(ctx context.Context, user entity.User) error
	GetProfile(ctx context.Context, email string) (interface{}, error)
	Logout(ctx context.Context, token string) error
	ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error)
}
