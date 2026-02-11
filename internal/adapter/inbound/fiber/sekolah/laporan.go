package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetReportData(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	req := model.ReportRequest{
		Type:   c.Query("type"),
		Period: c.Query("period"),
	}

	data, err := h.service.GetReportData(c.Context(), tenantID, req)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil data laporan", err)
	}
	return SendSuccess(c, "Data laporan berhasil diambil", data)
}
