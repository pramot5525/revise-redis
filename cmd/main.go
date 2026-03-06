package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/revise-redis/config"
	infraCache "github.com/revise-redis/infrastructure/cache"
	infraDB "github.com/revise-redis/infrastructure/db"
	"github.com/revise-redis/internal/adapters/http/handler"
	"github.com/revise-redis/internal/adapters/http/router"
	pgadapter "github.com/revise-redis/internal/adapters/postgres"
	redisadapter "github.com/revise-redis/internal/adapters/redis"
	"github.com/revise-redis/internal/app"
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

	// Secondary adapters (driven — implement output ports)
	newsRepo := pgadapter.NewNewsRepository(gormDB)
	newsCache := redisadapter.NewNewsCache(redisClient)

	// Use case (application layer)
	newsSvc := app.NewNewsService(newsRepo, newsCache, cfg.RedisTTL)

	// Primary adapters (driving — implement HTTP interface)
	newsHandler := handler.NewNewsHandler(newsSvc)

	// Fiber app
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	router.Register(app, newsHandler)

	log.Printf("server starting on :%s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
