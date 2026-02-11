package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetTahfidzSetoranList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetTahfidzSetoranList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar setoran", err)
	}
	return SendSuccess(c, "Daftar setoran berhasil diambil", data)
}

func (h *akademikHandler) CreateTahfidzSetoran(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, UpdatedAt, joined fields)
	var input struct {
		SantriID  string  `json:"santri_id"`
		UstadzID  *string `json:"ustadz_id"`
		Tanggal   string  `json:"tanggal"`
		Juz       int     `json:"juz"`
		Surah     string  `json:"surah"`
		AyatAwal  int     `json:"ayat_awal"`
		AyatAkhir int     `json:"ayat_akhir"`
		Tipe      string  `json:"tipe"`
		Kualitas  string  `json:"kualitas"`
		Catatan   string  `json:"catatan"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	m := model.TahfidzSetoran{
		TenantID:  tenantID, // From JWT, not user input
		SantriID:  input.SantriID,
		UstadzID:  input.UstadzID,
		Juz:       input.Juz,
		Surah:     input.Surah,
		AyatAwal:  input.AyatAwal,
		AyatAkhir: input.AyatAkhir,
		Tipe:      input.Tipe,
		Kualitas:  input.Kualitas,
		Catatan:   input.Catatan,
	}

	if err := h.service.CreateTahfidzSetoran(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mencatat setoran", err)
	}
	return SendCreated(c, "Setoran recorded", m)
}
