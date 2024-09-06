package service

import (
	"DoAn/pkg/account"
	"DoAn/pkg/account/auth"
	"DoAn/pkg/account/entity"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
)

type UserServiceStruct struct {
	repository account.UserRepository
	authenSvc  auth.AuthenService
	logger     log.Logger
}

func (u UserServiceStruct) ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error) {
	data, err := u.repository.ChangePassFirstTime(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u UserServiceStruct) Login(ctx context.Context, user entity.Account) (*string, *string, error) {
	loginSuccess, err := u.repository.Login(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	if !loginSuccess {
		return nil, nil, errors.New("username or password incorrect")
	}
	accessToken, err := u.authenSvc.CreateAccessToken(ctx, user.Username)
	if err != nil {
		fmt.Printf("error while creating access token: %v", err)
		return nil, nil, err
	}
	refreshToken, _, err := u.authenSvc.CreateRefreshToken(ctx, user.Username, "")
	if err != nil {
		fmt.Printf("error while creating refresh token: %v", err)
		return nil, nil, err
	}

	return &accessToken, &refreshToken, nil
}

func (u UserServiceStruct) Register(ctx context.Context, user entity.Account) (string, error) {
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

func (u UserServiceStruct) CheckExistedBarber(ctx context.Context, id int) (interface{}, error) {
	data, err := u.repository.CheckExistedBarber(ctx, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//func (u UserServiceStruct) RefreshToken(ctx context.Context, username string) (interface{}, error) {
//	data, err := u.repository.RefreshToken(ctx, username)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}

func (u UserServiceStruct) Logout(ctx context.Context, token string) (string, error) {
	var msg = "success"
	err := u.repository.Logout(ctx, token)
	if err != nil {
		return "error while logging out", err
	}
	return msg, nil
}

func NewService(rep account.UserRepository, svc auth.AuthenService, logger log.Logger) account.UserService {
	return &UserServiceStruct{
		repository: rep,
		logger:     logger,
		authenSvc:  svc,
	}
}
