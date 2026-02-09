package fiber_inbound_adapter

import (
	"context"

	"prabogo/internal/domain/landing_content"

	"github.com/gofiber/fiber/v2"
)

type landingContentHandler struct {
	domain *landing_content.Domain
}

func NewLandingContentHandler(domain *landing_content.Domain) *landingContentHandler {
	return &landingContentHandler{domain: domain}
}

// GET /api/public/landing/:key
func (h *landingContentHandler) Get(c *fiber.Ctx) error {
	ctx := context.Background()
	key := c.Params("key")

	content, err := h.domain.Get(ctx, key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	if content == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "Content not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   content,
	})
}

// PUT /api/v1/owner/landing/:key
func (h *landingContentHandler) Set(c *fiber.Ctx) error {
	ctx := context.Background()
	key := c.Params("key")

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	if err := h.domain.Set(ctx, key, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Content updated successfully",
	})
}
