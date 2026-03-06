package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/revise-redis/internal/core/domain"
	"github.com/revise-redis/internal/core/ports/input"
)

type NewsHandler struct {
	svc input.NewsService
}

func NewNewsHandler(svc input.NewsService) *NewsHandler {
	return &NewsHandler{svc: svc}
}

func (h *NewsHandler) GetAll(c *fiber.Ctx) error {
	news, err := h.svc.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    news,
		"message": "news fetched successfully",
	})
}

func (h *NewsHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}
	news, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "news not found",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    news,
		"message": "news fetched successfully",
	})
}

func (h *NewsHandler) Create(c *fiber.Ctx) error {
	var news domain.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}
	if err := h.svc.Create(&news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    news,
		"message": "news created successfully",
	})
}

func (h *NewsHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}
	var input domain.News
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}
	if err := h.svc.Update(id, &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "news updated successfully",
	})
}

func (h *NewsHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "news deleted successfully",
	})
}

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	return uint(id), err
}
