package mongodb

import (
	"DoAn/database"
	"DoAn/entity"
	"DoAn/pkg/user"
	"DoAn/pkg/user/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type repo struct {
	collection *mongo.Collection
	db         *sql.DB
	logger     log.Logger
}

func (r repo) ChangePassFirstTime(ctx context.Context, username, password string) (interface{}, error) {
	var result entity.User
	err := r.db.QueryRow("SELECT id,username,email,password,roles FROM userinfo WHERE username=$1 and password ='' ", username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Roles)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while selecting user info: %v", err))
		return nil, err
	}
	password, err = entity.Hash(password)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while hashing password: %v", err))
		return nil, err
	}
	_, err = r.db.Exec("UPDATE userinfo SET password=$1 WHERE username=$2", password, username)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while updating data: %v", err))
		return nil, err
	}
	return "success", nil
}

func (r repo) Login(ctx context.Context, user entity.User) (string, error) {
	username := entity.Santize(user.Username)
	password := entity.Santize(user.Password)
	var result entity.User
	err := r.db.QueryRow("SELECT id,username,email,password,roles FROM userinfo WHERE username=$1", username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Roles)
	if err != nil {
		return "error while selecting user info", errors.New("username not found")
	}

	hashedPassword := fmt.Sprintf("%v", result.Password)
	err = entity.CheckPasswordHash(hashedPassword, password)

	if err != nil {
		return "error generate jwt", errors.New("username or password incorrect")
	}
	token, errCreate := auth.Create(user.Username)
	if errCreate != nil {
		return "error generate jwt", errCreate
	}
	newToken := bson.M{"token": token, "user": user.Username, "created_at": time.Now()}
	_, errs := r.collection.InsertOne(context.TODO(), newToken)
	if errs != nil {
		return "error generate jwt", errs
	}
	return token, nil
}

func (r repo) Register(ctx context.Context, user entity.User) error {
	createTable := `
		CREATE TABLE IF NOT EXISTS userinfo (
			id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL ,
            email VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
		    roles VARCHAR(255) NOT NULL 
		);
	`
	_, err := r.db.Exec(createTable)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while creating table: %v", err))
		return err
	}
	var result entity.User
	errFindUsername := r.db.QueryRow("SELECT id,username,email,password,roles FROM userinfo WHERE username=$1", user.Username).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Roles)
	errFindEmail := r.db.QueryRow("SELECT id,username,email,password,roles FROM userinfo WHERE email=$1", user.Email).Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.Roles)

	if errFindUsername == nil || errFindEmail == nil {
		return errors.New("username or email is already exist")
	}
	password, err := entity.Hash(user.Password)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while hashing password: %v", err))
		return err
	}

	newUser := entity.User{Username: user.Username, Email: user.Email, Password: password, Roles: user.Roles}
	_, err = r.db.Exec("INSERT INTO userinfo (username, email, password,roles) VALUES ($1, $2, $3,$4)", newUser.Username, newUser.Email, newUser.Password, newUser.Roles)
	if err != nil {
		r.logger.Log(fmt.Sprintf("error while inserting data: %v", err))
		return err
	}

	return nil
}

func (r repo) GetProfile(ctx context.Context, email string) (interface{}, error) {
	var result entity.User
	err := r.db.QueryRow("SELECT id,username,email,roles FROM userinfo WHERE email=$1", email).Scan(&result.ID, &result.Username, &result.Email, &result.Roles)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r repo) Logout(ctx context.Context, tokenString string) error {
	collection := database.ConnectMongo(os.Getenv("TokenBlackListCollectionMongo"))
	_, err := collection.InsertOne(context.TODO(), bson.M{
		"token":          tokenString,
		"blacklisted_at": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(collection *mongo.Collection, collectionPostgres *sql.DB, logger log.Logger) (user.UserRepository, error) {
	return &repo{
		collection: collection,
		db:         collectionPostgres,
		logger:     logger,
	}, nil
}
