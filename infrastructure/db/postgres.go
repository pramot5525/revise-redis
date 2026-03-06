package db

import (
	"github.com/revise-redis/config"
	gpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(gpostgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
