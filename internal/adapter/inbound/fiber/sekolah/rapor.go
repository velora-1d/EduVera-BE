package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetRaporList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetRaporList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *akademikHandler) CreateRapor(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var m model.Rapor
	if err := c.BodyParser(&m); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Default status if not provided
	if m.Status == "" {
		m.Status = "Draft"
	}

	if err := h.service.CreateRapor(c.Context(), tenantID, &m); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Rapor created", "data": m})
}
