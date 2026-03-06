package postgres

import (
	"github.com/revise-redis/internal/core/domain"
	"gorm.io/gorm"
)

// newsModel is the GORM model — kept separate from the domain entity.
type newsModel struct {
	gorm.Model
	Title   string
	Content string
	Author  string
}

func (newsModel) TableName() string { return "news" }

func toModel(n *domain.News) *newsModel {
	return &newsModel{
		Model:   gorm.Model{ID: n.ID},
		Title:   n.Title,
		Content: n.Content,
		Author:  n.Author,
	}
}

func toDomain(m *newsModel) *domain.News {
	return &domain.News{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		Author:    m.Author,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// NewsRepository is the Postgres secondary adapter.
type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) FindAll() ([]domain.News, error) {
	var models []newsModel
	if err := r.db.Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	news := make([]domain.News, len(models))
	for i, m := range models {
		m := m
		news[i] = *toDomain(&m)
	}
	return news, nil
}

func (r *NewsRepository) FindByID(id uint) (*domain.News, error) {
	var m newsModel
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *NewsRepository) Create(n *domain.News) error {
	m := toModel(n)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	n.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *NewsRepository) Update(n *domain.News) error {
	m := toModel(n)
	return r.db.Save(m).Error
}

func (r *NewsRepository) Delete(id uint) error {
	return r.db.Delete(&newsModel{}, id).Error
}
