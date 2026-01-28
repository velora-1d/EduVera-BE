package pesantren

import (
	"eduvera/internal/domain/pesantren/dashboard"

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
	// TODO: Get tenant_id from context/token after middleware implementation
	// For now, use a dummy or query param if needed, or just pass a placeholder
	tenantID := "default-tenant"

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
