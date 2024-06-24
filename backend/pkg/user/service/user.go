package service

import (
	"DoAn/entity"
	"DoAn/pkg/user"
	"context"
	"github.com/go-kit/kit/log"
)

type UserServiceStruct struct {
	repository user.UserRepository
	logger     log.Logger
}

func (u UserServiceStruct) ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error) {
	data, err := u.repository.ChangePassFirstTime(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u UserServiceStruct) Login(ctx context.Context, user entity.User) (string, error) {
	jwt, err := u.repository.Login(ctx, user)
	if err != nil {
		return "error while generating jwt", err
	}
	return jwt, nil
}

func (u UserServiceStruct) Register(ctx context.Context, user entity.User) (string, error) {
	var msg = "success"
	if err := u.repository.Register(ctx, user); err != nil {
		errMsg := err.Error()

		return errMsg, err
	}
	return msg, nil
}

func (u UserServiceStruct) GetProfile(ctx context.Context, email string) (interface{}, error) {
	data, err := u.repository.GetProfile(ctx, email)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u UserServiceStruct) Logout(ctx context.Context, token string) (string, error) {
	var msg = "success"
	err := u.repository.Logout(ctx, token)
	if err != nil {
		return "error while logging out", err
	}
	return msg, nil
}

func NewService(rep user.UserRepository, logger log.Logger) user.UserService {
	return &UserServiceStruct{
		repository: rep,
		logger:     logger,
	}
}
