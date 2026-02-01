package fiber_inbound_adapter

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
)

type analyticsAdapter struct {
	domain domain.Domain
}

func NewAnalyticsAdapter(domain domain.Domain) *analyticsAdapter {
	return &analyticsAdapter{domain: domain}
}

// GET /api/v1/analytics
func (h *analyticsAdapter) GetAnalytics(c *fiber.Ctx) error {
	ctx := context.Background()

	tenantID := c.Locals("tenantID")
	if tenantID == nil || tenantID.(string) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID tidak ditemukan. Silakan login kembali.",
		})
	}

	analytics, err := h.domain.Analytics().GetAnalytics(ctx, tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data analytics.",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   analytics,
	})
}
