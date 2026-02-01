package fiber_inbound_adapter

import (
	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain/export"
	inbound_port "prabogo/internal/port/inbound"
)

type exportHandler struct {
	domain export.ExportDomain
}

func NewExportHandler(domain export.ExportDomain) inbound_port.ExportHttpPort {
	return &exportHandler{domain: domain}
}

// ExportStudents handles GET /api/v1/export/students?format=pdf|xlsx
func (h *exportHandler) ExportStudents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	format := c.Query("format", "pdf")

	if format != "pdf" && format != "xlsx" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format harus 'pdf' atau 'xlsx'",
		})
	}

	data, filename, err := h.domain.ExportStudents(c.Context(), tenantID, format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal export data siswa: " + err.Error(),
		})
	}

	contentType := "application/pdf"
	if format == "xlsx" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(data)
}

// ExportPayments handles GET /api/v1/export/payments?format=pdf|xlsx
func (h *exportHandler) ExportPayments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	format := c.Query("format", "pdf")

	if format != "pdf" && format != "xlsx" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format harus 'pdf' atau 'xlsx'",
		})
	}

	data, filename, err := h.domain.ExportPayments(c.Context(), tenantID, format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal export data keuangan: " + err.Error(),
		})
	}

	contentType := "application/pdf"
	if format == "xlsx" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(data)
}
