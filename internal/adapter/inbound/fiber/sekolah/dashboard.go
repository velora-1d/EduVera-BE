package sekolah

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetDashboardStats(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.service.GetDashboardStats(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": stats})
}
