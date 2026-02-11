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
		if err != nil {
			return SendError(c, http.StatusInternalServerError, "Gagal mengambil data rapor", err)
		}
	}
	return SendSuccess(c, "Data rapor berhasil diambil", data)
}

func (h *akademikHandler) CreateRapor(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, UpdatedAt, joined fields)
	var input struct {
		PeriodeID        string `json:"periode_id"`
		SantriID         string `json:"santri_id"`
		Status           string `json:"status"`
		CatatanWaliKelas string `json:"catatan_wali_kelas"`
	}
	if err := c.BodyParser(&input); err != nil {
		if err := c.BodyParser(&input); err != nil {
			return SendError(c, http.StatusBadRequest, "Invalid request body", err)
		}
	}

	// Default status if not provided
	status := input.Status
	if status == "" {
		status = "Draft"
	}

	// Explicit mapping: DTO → DB Model
	m := model.Rapor{
		TenantID:         tenantID, // From JWT, not user input
		PeriodeID:        input.PeriodeID,
		SantriID:         input.SantriID,
		Status:           status,
		CatatanWaliKelas: input.CatatanWaliKelas,
	}

	if err := h.service.CreateRapor(c.Context(), tenantID, &m); err != nil {
		if err := h.service.CreateRapor(c.Context(), tenantID, &m); err != nil {
			return SendError(c, http.StatusInternalServerError, "Gagal membuat rapor", err)
		}
	}
	return SendCreated(c, "Rapor created", m)
}
