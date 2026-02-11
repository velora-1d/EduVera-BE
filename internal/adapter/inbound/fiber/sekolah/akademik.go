package sekolah

import (
	"net/http"
	"prabogo/internal/domain/sekolah"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"

	"github.com/gofiber/fiber/v2"
)

type akademikHandler struct {
	service sekolah.AkademikDomain
}

func NewAkademikHandler(service sekolah.AkademikDomain) inbound_port.SekolahHttpPort {
	return &akademikHandler{
		service: service,
	}
}

// ------ Siswa Handler ------

func (h *akademikHandler) GetSiswaList(c *fiber.Ctx) error {
	// Get TenantID from context (set by middleware)
	tenantID := c.Locals("tenant_id").(string)

	siswaList, err := h.service.GetSiswaList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar siswa", err)
	}

	return SendSuccess(c, "Daftar siswa berhasil diambil", siswaList)
}

func (h *akademikHandler) CreateSiswa(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, etc.)
	var input struct {
		NIS      string `json:"nis"`
		Nama     string `json:"nama"`
		KelasID  string `json:"kelas_id"`
		Alamat   string `json:"alamat"`
		NamaWali string `json:"nama_wali"`
		NoHPWali string `json:"no_hp_wali"`
		Status   string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	siswa := model.Siswa{
		TenantID: tenantID, // From JWT, not user input
		NIS:      input.NIS,
		Nama:     input.Nama,
		KelasID:  input.KelasID,
		Alamat:   input.Alamat,
		NamaWali: input.NamaWali,
		NoHPWali: input.NoHPWali,
		Status:   input.Status,
	}

	if err := h.service.CreateSiswa(c.Context(), siswa); err != nil {
		return SendError(c, http.StatusInternalServerError, "Gagal membuat siswa", err)
	}

	return SendCreated(c, "Siswa created successfully", nil)
}

// ------ Guru Handler ------

func (h *akademikHandler) GetGuruList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	guruList, err := h.service.GetGuruList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return SendSuccess(c, "Daftar guru berhasil diambil", guruList)
}

func (h *akademikHandler) CreateGuru(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields
	var input struct {
		NIP    string `json:"nip"`
		Nama   string `json:"nama"`
		Jenis  string `json:"jenis"`
		Status string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	guru := model.Guru{
		TenantID: tenantID, // From JWT, not user input
		NIP:      input.NIP,
		Nama:     input.Nama,
		Jenis:    input.Jenis,
		Status:   input.Status,
	}

	if err := h.service.CreateGuru(c.Context(), guru); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat guru", err)
	}

	return SendCreated(c, "Guru created successfully", nil)
}

// ------ Mapel Handler ------

func (h *akademikHandler) GetMapelList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	mapelList, err := h.service.GetMapelList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar mapel", err)
	}

	return SendSuccess(c, "Daftar mapel berhasil diambil", mapelList)
}

// ------ Kelas Handler ------

func (h *akademikHandler) GetKelasList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	kelasList, err := h.service.GetKelasList(c.Context(), tenantID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar kelas", err)
	}

	return SendSuccess(c, "Daftar kelas berhasil diambil", kelasList)
}

func (h *akademikHandler) CreateKelas(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields
	var input struct {
		Nama    string `json:"nama"`
		Tingkat string `json:"tingkat"`
		Urutan  int    `json:"urutan"`
		Status  string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Explicit mapping: DTO → DB Model
	kelas := model.Kelas{
		TenantID: tenantID, // From JWT, not user input
		Nama:     input.Nama,
		Tingkat:  input.Tingkat,
		Urutan:   input.Urutan,
		Status:   input.Status,
	}

	if err := h.service.CreateKelas(c.Context(), tenantID, &kelas); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat kelas", err)
	}
	return SendCreated(c, "Kelas created", kelas)
}
