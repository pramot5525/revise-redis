package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/revise-redis/internal/core/domain"
)

const (
	keyAll    = "news:all"
	keyPrefix = "news:"
)

// NewsCache is the Redis secondary adapter.
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
	if err := json.Unmarshal(data, &news); err != nil {
		return nil, err
	}
	return news, nil
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
	if err := json.Unmarshal(data, &news); err != nil {
		return nil, err
	}
	return &news, nil
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
