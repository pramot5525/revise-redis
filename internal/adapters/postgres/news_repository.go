package postgres

import (
	"github.com/revise-redis/internal/domain"
	"gorm.io/gorm"
)

// newsModel is the internal GORM model — never exposed outside this package.
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

// NewsRepository implements port.NewsRepository using GORM + PostgreSQL.
type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) FindAll() ([]domain.News, error) {
	var rows []newsModel
	if err := r.db.Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	news := make([]domain.News, len(rows))
	for i, m := range rows {
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
	return r.db.Save(toModel(n)).Error
}

func (r *NewsRepository) Delete(id uint) error {
	return r.db.Delete(&newsModel{}, id).Error
}
