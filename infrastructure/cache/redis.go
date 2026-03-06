package cache

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
	"github.com/revise-redis/config"
)

func NewRedis(cfg *config.Config) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
