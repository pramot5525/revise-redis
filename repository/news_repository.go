package repository

import (
	"github.com/revise-redis/models"
	"gorm.io/gorm"
)

type NewsRepository interface {
	FindAll() ([]models.News, error)
	FindByID(id uint) (*models.News, error)
	Create(news *models.News) error
	Update(news *models.News) error
	Delete(id uint) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) FindAll() ([]models.News, error) {
	var news []models.News
	if err := r.db.Order("created_at desc").Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsRepository) FindByID(id uint) (*models.News, error) {
	var news models.News
	if err := r.db.First(&news, id).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) Create(news *models.News) error {
	return r.db.Create(news).Error
}

func (r *newsRepository) Update(news *models.News) error {
	return r.db.Save(news).Error
}

func (r *newsRepository) Delete(id uint) error {
	return r.db.Delete(&models.News{}, id).Error
}
