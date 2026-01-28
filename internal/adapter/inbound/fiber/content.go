package fiber_inbound_adapter

import (
	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"

	"github.com/gofiber/fiber/v2"
)

type contentAdapter struct {
	domain domain.Domain
}

func NewContentAdapter(domain domain.Domain) inbound_port.ContentHttpPort {
	return &contentAdapter{
		domain: domain,
	}
}

func (h *contentAdapter) Get(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Key is required",
		})
	}

	content, err := h.domain.Content().Get(c.Context(), key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	return c.JSON(content)
}

func (h *contentAdapter) Upsert(c *fiber.Ctx) error {
	var input model.ContentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	content, err := h.domain.Content().Upsert(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save content",
		})
	}

	return c.JSON(content)
}
