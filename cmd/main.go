package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/revise-redis/cache"
	"github.com/revise-redis/config"
	"github.com/revise-redis/db"
	"github.com/revise-redis/handler"
	"github.com/revise-redis/models"
	"github.com/revise-redis/repository"
	"github.com/revise-redis/router"
	"github.com/revise-redis/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database
	gormDB, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	if err := gormDB.AutoMigrate(&models.News{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Redis
	rdb, err := cache.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Layers
	newsRepo := repository.NewNewsRepository(gormDB)
	newsSvc := service.NewNewsService(newsRepo, rdb, cfg.RedisTTL)
	newsHandler := handler.NewNewsHandler(newsSvc)

	// Fiber app
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	router.Register(app, newsHandler)

	log.Printf("server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
