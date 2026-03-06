package output

import (
	"time"

	"github.com/revise-redis/internal/core/domain"
)

// NewsCache is the secondary (driven) port for caching.
type NewsCache interface {
	GetAll() ([]domain.News, error)
	SetAll(news []domain.News, ttl time.Duration) error
	GetByID(id uint) (*domain.News, error)
	SetByID(news *domain.News, ttl time.Duration) error
	DeleteByID(id uint) error
	DeleteAll() error
}
