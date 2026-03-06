package input

import "github.com/revise-redis/internal/core/domain"

// NewsService is the primary (driving) port — defines what the application can do.
type NewsService interface {
	GetAll() ([]domain.News, error)
	GetByID(id uint) (*domain.News, error)
	Create(news *domain.News) error
	Update(id uint, news *domain.News) error
	Delete(id uint) error
}
