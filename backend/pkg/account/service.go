package account

import (
	"DoAn/entity"
	"context"
)

type UserService interface {
	Login(ctx context.Context, user entity.Account) (*string, *string, error)
	Register(ctx context.Context, user entity.Account) (string, error)

	GetProfile(ctx context.Context, email string) (interface{}, error)
	Logout(ctx context.Context, token string) (string, error)
	ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error)
	//RefreshToken(ctx context.Context, username string) (interface{}, error)

	CheckExistedBarber(ctx context.Context, id int) (interface{}, error)

	GetAllBarber(ctx context.Context) (interface{}, error)
}
