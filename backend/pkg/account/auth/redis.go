package auth

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(address, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to ping Redis cache")
		return nil
	}
	fmt.Println(pong + " Redis cache connected")

	return rdb
}
