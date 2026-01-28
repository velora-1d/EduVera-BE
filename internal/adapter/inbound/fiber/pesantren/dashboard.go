package pesantren

import (
	"prabogo/internal/domain/pesantren/dashboard"

	"github.com/gofiber/fiber/v2"
)

type DashboardHttpPort interface {
	GetStats(c *fiber.Ctx) error
}

type dashboardAdapter struct {
	service dashboard.Service
}

func NewDashboardAdapter(service dashboard.Service) DashboardHttpPort {
	return &dashboardAdapter{
		service: service,
	}
}

func (h *dashboardAdapter) GetStats(c *fiber.Ctx) error {
	// Get tenant_id from context (set by middleware)
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.service.GetStats(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch dashboard stats",
		})
	}

	return c.JSON(fiber.Map{
		"data": stats,
	})
}
