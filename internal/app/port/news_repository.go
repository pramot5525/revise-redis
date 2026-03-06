package port

import "github.com/revise-redis/internal/domain"

// NewsRepository is the output port for persistence.
type NewsRepository interface {
	FindAll() ([]domain.News, error)
	FindByID(id uint) (*domain.News, error)
	Create(news *domain.News) error
	Update(news *domain.News) error
	Delete(id uint) error
}
