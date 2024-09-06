package account

import (
	"DoAn/pkg/account/entity"
	"context"
)

type UserRepository interface {
	Register(ctx context.Context, user entity.Account) error
	GetProfile(ctx context.Context, email string) (interface{}, error)

	Login(ctx context.Context, user entity.Account) (bool, error)
	Logout(ctx context.Context, token string) error
	ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error)
	//RefreshToken(ctx context.Context, username string) (interface{}, error)

	CheckExistedBarber(ctx context.Context, id int) (string, error)
}
