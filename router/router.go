package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revise-redis/handler"
)

func Register(app *fiber.App, h *handler.NewsHandler) {
	api := app.Group("/api")

	news := api.Group("/news")
	news.Get("/", h.GetAll)
	news.Get("/:id", h.GetByID)
	news.Post("/", h.Create)
	news.Put("/:id", h.Update)
	news.Delete("/:id", h.Delete)
}
