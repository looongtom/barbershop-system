package db

import (
	"DoAn/pkg/account"
	"DoAn/pkg/account/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	collection         *mongo.Collection
	blkTokenCollection *mongo.Collection
	db                 *sql.DB
	logger             log.Logger
}

func (r repo) ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error) {
	var result entity.Account
	err := r.db.QueryRow("SELECT id,username,email,password,Role FROM account WHERE username=$1 and password ='' ", username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Role)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while selecting account info: %v", err))
		return nil, err
	}
	password, err = entity.Hash(password)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while hashing password: %v", err))
		return nil, err
	}
	_, err = r.db.Exec("UPDATE account SET password=$1 WHERE username=$2", password, username)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while updating data: %v", err))
		return nil, err
	}
	return "success", nil
}

func (r repo) Login(ctx context.Context, user entity.Account) (*int, error) {
	username := entity.Santize(user.Username)
	password := entity.Santize(user.Password)
	var result entity.Account
	err := r.db.QueryRow("SELECT id,username,email,password,role FROM account WHERE username=$1", username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Role)
	if err != nil {
		return nil, errors.New("username not found")
	}

	hashedPassword := fmt.Sprintf("%v", result.Password)
	err = entity.CheckPasswordHash(hashedPassword, password)

	if err != nil {
		return nil, errors.New("username or password incorrect")
	}

	//token, errCreate := auth.CreateAccessToken(user.Username)
	//if errCreate != nil {
	//	return false, errCreate
	//}
	//
	//newToken := bson.M{"token": token, "account": user.Username, "created_at": time.Now()}
	//_, errs := r.collection.InsertOne(context.TODO(), newToken)
	//if errs != nil {
	//	return nil, nil, errs
	//}
	//
	//refreshToken, errCreate := auth.CreateRefreshToken(user.Username)
	//if errCreate != nil {
	//	return nil, nil, errCreate
	//}
	//newRefreshToken := bson.M{"refresh_token": refreshToken, "account": user.Username, "created_at": time.Now()}
	//_, errs = r.collection.InsertOne(context.TODO(), newRefreshToken)
	//if errs != nil {
	//	return nil, nil, errs
	//}

	return &result.ID, nil
}

func (r repo) Register(ctx context.Context, user entity.Account) error {
	createTable := `
		CREATE TABLE IF NOT EXISTS account (
			id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL ,
            email VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
		    role int NOT NULL ,
		    phone_number VARCHAR(255) NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			about VARCHAR(255) NOT NULL,
			avatar VARCHAR(255) NOT NULL,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
	`
	_, err := r.db.Exec(createTable)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while creating table: %v", err))
		return err
	}
	var result entity.Account
	errFindUsername := r.db.QueryRow("SELECT id,username,email,password,role FROM account WHERE username=$1", user.Username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Role)
	errFindEmail := r.db.QueryRow("SELECT id,username,email,password,role FROM account WHERE email=$1", user.Email).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Role)

	if errFindUsername == nil || errFindEmail == nil {
		return errors.New("username or email is already exist")
	}
	password, err := entity.Hash(user.Password)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while hashing password: %v", err))
		return err
	}

	newUser := entity.Account{
		Username:    user.Username,
		Email:       user.Email,
		Password:    password,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
		About:       user.About,
		Avatar:      user.Avatar,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	_, err = r.db.Exec("INSERT INTO account (username, email, password,role,phone_number,full_name,about,avatar,created_at,updated_at) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9,$10)",
		newUser.Username, newUser.Email, newUser.Password, newUser.Role, newUser.PhoneNumber, newUser.FullName, newUser.About, newUser.Avatar, newUser.CreatedAt, newUser.UpdatedAt)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while inserting data: %v", err))
		return err
	}

	return nil
}

//func (r repo) RefreshToken(ctx context.Context, username string) (interface{}, error) {
//	var result entity.Account
//	err := r.db.QueryRow("SELECT id,username,email,password,role FROM account WHERE username=$1", username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Role)
//	if err != nil {
//		return nil, errors.New("username not found")
//	}
//	refreshToken, errCreate := auth.CreateRefreshToken(username)
//	if errCreate != nil {
//		return nil, errCreate
//	}
//	newRefreshToken := bson.M{"refresh_token": refreshToken, "account": username, "created_at": time.Now()}
//	_, errs := r.collection.InsertOne(context.TODO(), newRefreshToken)
//	if errs != nil {
//		return nil, errs
//	}
//	return refreshToken, nil
//}

func (r repo) CheckExistedBarber(ctx context.Context, id int) (string, error) {
	query := `
	SELECT id,role from account where id = $1 ;
`
	checkedId, checkedRole := -1, -1
	err := r.db.QueryRow(query, id).Scan(&checkedId, &checkedRole)
	roleName := "UNKNOWN"
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return roleName, err
	}
	switch checkedRole {
	case entity.RoleBarber:
		roleName = "BARBER"
	case entity.RoleUser:
		roleName = "USER"
	case entity.RoleAdmin:
		roleName = "ADMIN"
	default:
		roleName = "UNKNOWN"
	}
	return roleName, nil
}
func (r repo) GetProfile(ctx context.Context, email string) (interface{}, error) {
	var result entity.Account
	err := r.db.QueryRow("SELECT id,username,email,role,phone_number,full_name,about,avatar,created_at,updated_at FROM account WHERE email=$1", email).
		Scan(&result.ID, &result.Username, &result.Email, &result.Role, &result.PhoneNumber, &result.FullName, &result.About, &result.Avatar, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r repo) Logout(ctx context.Context, tokenString string) error {
	_, err := r.blkTokenCollection.InsertOne(context.TODO(), bson.M{
		"token":          tokenString,
		"blacklisted_at": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(collection *mongo.Collection, blkListCollection *mongo.Collection, collectionPostgres *sql.DB, logger log.Logger) (account.UserRepository, error) {
	return &repo{
		collection:         collection,
		blkTokenCollection: blkListCollection,
		db:                 collectionPostgres,
		logger:             logger,
	}, nil
}
