package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetTabunganList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetTabunganList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar tabungan", err)
	}
	return SendSuccess(c, "Daftar tabungan berhasil diambil", data)
}

func (h *akademikHandler) CreateTabunganMutasi(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt)
	var input struct {
		TabunganID string `json:"tabungan_id"`
		Tipe       string `json:"tipe"`
		Nominal    int64  `json:"nominal"`
		Keterangan string `json:"keterangan"`
		Petugas    string `json:"petugas"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	m := model.TabunganMutasi{
		TenantID:   tenantID, // From JWT, not user input
		TabunganID: input.TabunganID,
		Tipe:       input.Tipe,
		Nominal:    input.Nominal,
		Keterangan: input.Keterangan,
		Petugas:    input.Petugas,
	}

	if err := h.service.CreateTabunganMutasi(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal membuat mutasi", err)
	}
	return SendCreated(c, "Mutasi created", m)
}
