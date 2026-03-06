package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/revise-redis/internal/domain"
)

const (
	keyAll    = "news:all"
	keyPrefix = "news:"
)

// NewsCache implements port.NewsCache using Redis.
type NewsCache struct {
	client *goredis.Client
}

func NewNewsCache(client *goredis.Client) *NewsCache {
	return &NewsCache{client: client}
}

func (c *NewsCache) GetAll() ([]domain.News, error) {
	data, err := c.client.Get(context.Background(), keyAll).Bytes()
	if err != nil {
		return nil, err
	}
	var news []domain.News
	return news, json.Unmarshal(data, &news)
}

func (c *NewsCache) SetAll(news []domain.News, ttl time.Duration) error {
	data, err := json.Marshal(news)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), keyAll, data, ttl).Err()
}

func (c *NewsCache) GetByID(id uint) (*domain.News, error) {
	data, err := c.client.Get(context.Background(), key(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var news domain.News
	return &news, json.Unmarshal(data, &news)
}

func (c *NewsCache) SetByID(news *domain.News, ttl time.Duration) error {
	data, err := json.Marshal(news)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), key(news.ID), data, ttl).Err()
}

func (c *NewsCache) DeleteByID(id uint) error {
	return c.client.Del(context.Background(), key(id)).Err()
}

func (c *NewsCache) DeleteAll() error {
	return c.client.Del(context.Background(), keyAll).Err()
}

func key(id uint) string {
	return fmt.Sprintf("%s%d", keyPrefix, id)
}
