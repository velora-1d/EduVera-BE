package sekolah

import (
	"net/http"
	"prabogo/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetPelanggaranAturanList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetPelanggaranAturanList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar pelanggaran aturan", err)
	}
	return SendSuccess(c, "Daftar pelanggaran aturan berhasil diambil", data)
}

func (h *akademikHandler) CreatePelanggaranAturan(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, UpdatedAt)
	var input struct {
		Judul    string `json:"judul"`
		Kategori string `json:"kategori"`
		Poin     int    `json:"poin"`
		Level    string `json:"level"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	m := model.PelanggaranAturan{
		TenantID: tenantID, // From JWT, not user input
		Judul:    input.Judul,
		Kategori: input.Kategori,
		Poin:     input.Poin,
		Level:    input.Level,
	}

	if err := h.service.CreatePelanggaranAturan(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal membuat aturan", err)
	}
	return SendCreated(c, "Aturan created", m)
}

func (h *akademikHandler) GetPelanggaranSiswaList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetPelanggaranSiswaList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar pelanggaran siswa", err)
	}
	return SendSuccess(c, "Daftar pelanggaran siswa berhasil diambil", data)
}

func (h *akademikHandler) CreatePelanggaranSiswa(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields
	var input struct {
		SantriID   string  `json:"santri_id"`
		AturanID   *string `json:"aturan_id"`
		Tanggal    string  `json:"tanggal"`
		Poin       int     `json:"poin"`
		Keterangan string  `json:"keterangan"`
		Status     string  `json:"status"`
		Sanksi     string  `json:"sanksi"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	m := model.PelanggaranSiswa{
		TenantID:   tenantID, // From JWT, not user input
		SantriID:   input.SantriID,
		AturanID:   input.AturanID,
		Poin:       input.Poin,
		Keterangan: input.Keterangan,
		Status:     input.Status,
		Sanksi:     input.Sanksi,
	}

	// Parse Tanggal
	if input.Tanggal != "" {
		if t, err := time.Parse("2006-01-02", input.Tanggal); err == nil {
			m.Tanggal = t
		}
	}

	if err := h.service.CreatePelanggaranSiswa(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mencatat pelanggaran", err)
	}
	return SendCreated(c, "Pelanggaran recorded", m)
}

func (h *akademikHandler) GetPerizinanList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetPerizinanList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar perizinan", err)
	}
	return SendSuccess(c, "Daftar perizinan berhasil diambil", data)
}

func (h *akademikHandler) CreatePerizinan(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields
	var input struct {
		SantriID string `json:"santri_id"`
		Tipe     string `json:"tipe"`
		Alasan   string `json:"alasan"`
		Dari     string `json:"dari"`
		Sampai   string `json:"sampai"`
		Status   string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, http.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	m := model.Perizinan{
		TenantID: tenantID, // From JWT, not user input
		SantriID: input.SantriID,
		Tipe:     input.Tipe,
		Alasan:   input.Alasan,
		Status:   input.Status,
	}

	// Parse date fields
	if input.Dari != "" {
		if t, err := time.Parse(time.RFC3339, input.Dari); err == nil {
			m.Dari = t
		}
	}
	if input.Sampai != "" {
		if t, err := time.Parse(time.RFC3339, input.Sampai); err == nil {
			m.Sampai = t
		}
	}

	if err := h.service.CreatePerizinan(c.Context(), tenantID, &m); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal membuat perizinan", err)
	}
	return SendCreated(c, "Perizinan created", m)
}
