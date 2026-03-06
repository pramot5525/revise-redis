package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/revise-redis/config"
	infraCache "github.com/revise-redis/infrastructure/cache"
	infraDB "github.com/revise-redis/infrastructure/db"
	"github.com/revise-redis/internal/adapters/primary/http/handler"
	"github.com/revise-redis/internal/adapters/primary/http/router"
	pgadapter "github.com/revise-redis/internal/adapters/secondary/postgres"
	redisadapter "github.com/revise-redis/internal/adapters/secondary/redis"
	"github.com/revise-redis/internal/core/service"
	"gorm.io/gorm"
)

// migrateModel is the GORM model used only for schema migration.
type migrateModel struct {
	gorm.Model
	Title   string
	Content string
	Author  string
}

func (migrateModel) TableName() string { return "news" }

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Infrastructure — connections
	gormDB, err := infraDB.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	if err := gormDB.AutoMigrate(&migrateModel{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	redisClient, err := infraCache.NewRedis(cfg)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	// Secondary adapters (driven)
	newsRepo := pgadapter.NewNewsRepository(gormDB)
	newsCache := redisadapter.NewNewsCache(redisClient)

	// Core service
	newsSvc := service.NewNewsService(newsRepo, newsCache, cfg.RedisTTL)

	// Primary adapters (driving)
	newsHandler := handler.NewNewsHandler(newsSvc)

	// Fiber app
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	router.Register(app, newsHandler)

	log.Printf("server starting on :%s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
