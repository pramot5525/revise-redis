package app

import (
	"time"

	"github.com/revise-redis/internal/app/port"
	"github.com/revise-redis/internal/domain"
)

type newsService struct {
	repo  port.NewsRepository
	cache port.NewsCache
	ttl   time.Duration
}

// NewNewsService returns an implementation of port.NewsService.
func NewNewsService(repo port.NewsRepository, cache port.NewsCache, ttlSeconds int) port.NewsService {
	return &newsService{
		repo:  repo,
		cache: cache,
		ttl:   time.Duration(ttlSeconds) * time.Second,
	}
}

func (s *newsService) GetAll() ([]domain.News, error) {
	if cached, err := s.cache.GetAll(); err == nil {
		return cached, nil
	}
	news, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetAll(news, s.ttl)
	return news, nil
}

func (s *newsService) GetByID(id uint) (*domain.News, error) {
	if cached, err := s.cache.GetByID(id); err == nil {
		return cached, nil
	}
	news, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetByID(news, s.ttl)
	return news, nil
}

func (s *newsService) Create(news *domain.News) error {
	if err := s.repo.Create(news); err != nil {
		return err
	}
	_ = s.cache.DeleteAll()
	return nil
}

func (s *newsService) Update(id uint, input *domain.News) error {
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
	_ = s.cache.DeleteByID(id)
	_ = s.cache.DeleteAll()
	return nil
}

func (s *newsService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	_ = s.cache.DeleteByID(id)
	_ = s.cache.DeleteAll()
	return nil
}
