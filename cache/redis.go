package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/revise-redis/config"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
