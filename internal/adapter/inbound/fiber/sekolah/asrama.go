package sekolah

import (
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

// Asrama Handlers
func (h *akademikHandler) GetAsramaList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	list, err := h.service.GetAsramaList(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": list})
}

func (h *akademikHandler) CreateAsrama(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req model.Asrama
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.service.CreateAsrama(c.Context(), tenantID, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Asrama created", "data": req})
}

// Kamar Handlers
func (h *akademikHandler) GetKamarList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	asramaID := c.Query("asrama_id")
	list, err := h.service.GetKamarList(c.Context(), tenantID, asramaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": list})
}

func (h *akademikHandler) CreateKamar(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req model.Kamar
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.service.CreateKamar(c.Context(), tenantID, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Kamar created", "data": req})
}

// Penempatan Handlers
func (h *akademikHandler) GetPenempatanList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	list, err := h.service.GetPenempatanList(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": list})
}

func (h *akademikHandler) CreatePenempatan(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req model.Penempatan
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.service.CreatePenempatan(c.Context(), tenantID, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Penempatan created", "data": req})
}
