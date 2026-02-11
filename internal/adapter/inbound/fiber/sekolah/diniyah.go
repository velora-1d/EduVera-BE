package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetDiniyahKitabList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetDiniyahKitabList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar kitab", err)
	}
	return SendSuccess(c, "Daftar kitab berhasil diambil", data)
}

func (h *akademikHandler) CreateDiniyahKitab(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, UpdatedAt)
	var input struct {
		NamaKitab   string `json:"nama_kitab"`
		BidangStudi string `json:"bidang_studi"`
		Pengarang   string `json:"pengarang"`
		Keterangan  string `json:"keterangan"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	m := model.DiniyahKitab{
		TenantID:    tenantID, // From JWT, not user input
		NamaKitab:   input.NamaKitab,
		BidangStudi: input.BidangStudi,
		Pengarang:   input.Pengarang,
		Keterangan:  input.Keterangan,
	}

	if err := h.service.CreateDiniyahKitab(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal membuat kitab", err)
	}
	return SendCreated(c, "Kitab created", m)
}
