package redis

import (
	"context"
	"fmt"
	"log"
	"msnserver/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(c config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.DB,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis server: %v", err)
	}
	if pong != "PONG" {
		return nil, fmt.Errorf("unexpected response from redis server: %s", pong)
	}

	log.Println("Redis initialized successfully")

	return rdb, nil
}
