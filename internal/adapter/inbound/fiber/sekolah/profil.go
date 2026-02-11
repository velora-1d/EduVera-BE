package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetProfil(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetProfil(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil data profil", err)
	}
	return SendSuccess(c, "Data profil berhasil diambil", data)
}

func (h *akademikHandler) UpdateProfil(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var m model.ProfilUpdate
	if err := c.BodyParser(&m); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.UpdateProfil(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengupdate profil", err)
	}
	return SendSuccess(c, "Profil updated", nil)
}
