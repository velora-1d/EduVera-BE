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
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil konten landing page", err)
	}

	if content == nil {
		return SendError(c, fiber.StatusNotFound, "Content not found", nil)
	}

	return SendSuccess(c, "Success", content)
}

// PUT /api/v1/owner/landing/:key
func (h *landingContentHandler) Set(c *fiber.Ctx) error {
	ctx := context.Background()
	key := c.Params("key")

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.domain.Set(ctx, key, input); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengupdate konten", err)
	}

	return SendSuccess(c, "Content updated successfully", nil)
}
