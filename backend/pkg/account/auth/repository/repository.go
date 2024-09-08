package repository

import (
	"DoAn/pkg/account/auth"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisTokenRepository struct {
	Redis *redis.Client
}

func (r RedisTokenRepository) GetAllTokenByUserid(ctx context.Context, userId string) (map[string]string, error) {
	//pattern := fmt.Sprintf("%s*", userId)
	results := make(map[string]string)
	iter := r.Redis.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		value, err := r.Redis.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get value for key %s: %v\n", key, err)
			continue
		}
		results[key] = value
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r RedisTokenRepository) RemoveTokenByUserid(ctx context.Context, userId string) error {
	result := r.Redis.Del(ctx, userId)
	if err := result.Err(); err != nil {
		log.Printf("Failed to delete key %s: %v\n", userId, err)
		return err
	}
	if result.Val() < 1 {
		log.Printf("Key %s not found\n", userId)
		return fmt.Errorf("key %s not found", userId)
	}
	return nil
}

func (r RedisTokenRepository) ValidateTokenInRedis(ctx context.Context, userId, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s", userId, token)
	result, err := r.Redis.Get(ctx, key).Int()
	if err != nil && err.Error() != "redis: nil" {
		log.Printf("Failed to get refresh token for userId/tokenId %s/%s: %v \n", userId, token, err)
		return false, err
	} else if err != nil && err.Error() == "redis: nil" {
		return true, nil
	}
	if result == 1 {
		return true, nil
	}
	return false, nil
}

func (r RedisTokenRepository) SetBlacklistToken(ctx context.Context, userId, token string) error {
	key := fmt.Sprintf("%s:%s", userId, token)
	if err := r.Redis.Set(ctx, key, 0, 0).Err(); err != nil {
		log.Printf("Failed to set refresh token for userId/tokenId %s/%s: %v \n", userId, token, err)
		return err
	}
	return nil
}

func (r RedisTokenRepository) SetRefreshToken(ctx context.Context,
	userId, token string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userId, token)
	if err := r.Redis.Set(ctx, key, 1, expiresIn).Err(); err != nil {
		log.Printf("Failed to set refresh token for userId/tokenId %s/%s: %v \n", userId, token, err)
		return err
	}
	return nil
}

func (r RedisTokenRepository) DeleteRefreshToken(ctx context.Context, userId, tokenId string) error {
	key := fmt.Sprintf("%s:%s", userId, tokenId)
	result := r.Redis.Del(ctx, key)
	if err := result.Err(); err != nil {
		log.Printf("Failed to delete refresh token for userId/tokenId %s/%s: %v \n", userId, tokenId, err)
		return err
	}

	if result.Val() < 1 {
		log.Printf("Refresh token for userId/tokenId %s/%s not found \n", userId, tokenId)
		return errors.New("Invalid refresh token")
	}
	return nil
}

func (r RedisTokenRepository) DeleteUserRefreshToken(ctx context.Context, userId string) error {
	pattern := fmt.Sprintf("%s*", userId)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		err := r.Redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Printf("Failed to delete refresh token for userId/tokenId %s: %v \n", iter.Val(), err)
			failCount++
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token for tokenId %s\n", iter.Val())
	}

	if failCount > 0 {
		return errors.New("Failed to delete user refresh token")
	}
	return nil
}

func NewTokenRepository(redisClient *redis.Client) auth.TokenRepository {
	return &RedisTokenRepository{
		Redis: redisClient,
	}
}
