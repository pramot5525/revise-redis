package port

import "github.com/revise-redis/internal/domain"

// NewsService is the input port — defines what the application exposes to the outside world.
type NewsService interface {
	GetAll() ([]domain.News, error)
	GetByID(id uint) (*domain.News, error)
	Create(news *domain.News) error
	Update(id uint, news *domain.News) error
	Delete(id uint) error
}
