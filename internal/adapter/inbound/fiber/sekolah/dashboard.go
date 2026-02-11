package sekolah

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetDashboardStats(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.service.GetDashboardStats(c.Context(), tenantID)
	if err != nil {
		if err != nil {
			return SendError(c, http.StatusInternalServerError, "Gagal mengambil statistik dashboard", err)
		}
	}
	return SendSuccess(c, "Statistik dashboard berhasil diambil", stats)
}
