package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/revise-redis/models"
	"github.com/revise-redis/repository"
)

const (
	keyAll    = "news:all"
	keyPrefix = "news:"
)

type NewsService interface {
	GetAll() ([]models.News, error)
	GetByID(id uint) (*models.News, error)
	Create(news *models.News) error
	Update(id uint, news *models.News) error
	Delete(id uint) error
}

type newsService struct {
	repo repository.NewsRepository
	rdb  *redis.Client
	ttl  time.Duration
}

func NewNewsService(repo repository.NewsRepository, rdb *redis.Client, ttlSeconds int) NewsService {
	return &newsService{
		repo: repo,
		rdb:  rdb,
		ttl:  time.Duration(ttlSeconds) * time.Second,
	}
}

func (s *newsService) GetAll() ([]models.News, error) {
	ctx := context.Background()

	cached, err := s.rdb.Get(ctx, keyAll).Result()
	if err == nil {
		var news []models.News
		if err := json.Unmarshal([]byte(cached), &news); err == nil {
			return news, nil
		}
	}

	news, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(news); err == nil {
		s.rdb.Set(ctx, keyAll, data, s.ttl)
	}

	return news, nil
}

func (s *newsService) GetByID(id uint) (*models.News, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", keyPrefix, id)

	cached, err := s.rdb.Get(ctx, key).Result()
	if err == nil {
		var news models.News
		if err := json.Unmarshal([]byte(cached), &news); err == nil {
			return &news, nil
		}
	}

	news, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(news); err == nil {
		s.rdb.Set(ctx, key, data, s.ttl)
	}

	return news, nil
}

func (s *newsService) Create(news *models.News) error {
	if err := s.repo.Create(news); err != nil {
		return err
	}
	s.invalidate()
	return nil
}

func (s *newsService) Update(id uint, input *models.News) error {
	news, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	news.Title = input.Title
	news.Content = input.Content
	news.Author = input.Author

	if err := s.repo.Update(news); err != nil {
		return err
	}

	ctx := context.Background()
	s.rdb.Del(ctx, fmt.Sprintf("%s%d", keyPrefix, id))
	s.rdb.Del(ctx, keyAll)
	return nil
}

func (s *newsService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	ctx := context.Background()
	s.rdb.Del(ctx, fmt.Sprintf("%s%d", keyPrefix, id))
	s.rdb.Del(ctx, keyAll)
	return nil
}

func (s *newsService) invalidate() {
	ctx := context.Background()
	s.rdb.Del(ctx, keyAll)
}
