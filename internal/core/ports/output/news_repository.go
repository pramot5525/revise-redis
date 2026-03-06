package output

import "github.com/revise-redis/internal/core/domain"

// NewsRepository is the secondary (driven) port for persistence.
type NewsRepository interface {
	FindAll() ([]domain.News, error)
	FindByID(id uint) (*domain.News, error)
	Create(news *domain.News) error
	Update(news *domain.News) error
	Delete(id uint) error
}
