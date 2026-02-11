package sekolah

import (
	"prabogo/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Asrama Handlers
func (h *akademikHandler) GetAsramaList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	list, err := h.service.GetAsramaList(c.Context(), tenantID)
	if err != nil {
		if err != nil {
			return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data asrama", err)
		}
	}
	return SendSuccess(c, "Data asrama berhasil diambil", list)
}

func (h *akademikHandler) CreateAsrama(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, calculated fields)
	var input struct {
		Nama      string  `json:"nama"`
		Jenis     string  `json:"jenis"`
		MusyrifID *string `json:"musyrif_id"`
		Status    string  `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		if err := c.BodyParser(&input); err != nil {
			return SendError(c, fiber.StatusBadRequest, "Invalid request", err)
		}
	}

	// Explicit mapping: DTO → DB Model
	req := model.Asrama{
		TenantID:  tenantID, // From JWT, not user input
		Nama:      input.Nama,
		Jenis:     input.Jenis,
		MusyrifID: input.MusyrifID,
		Status:    input.Status,
	}

	if err := h.service.CreateAsrama(c.Context(), tenantID, &req); err != nil {
		if err := h.service.CreateAsrama(c.Context(), tenantID, &req); err != nil {
			return SendError(c, fiber.StatusInternalServerError, "Gagal membuat asrama", err)
		}
	}
	return SendCreated(c, "Asrama created", req)
}

// Kamar Handlers
func (h *akademikHandler) GetKamarList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	asramaID := c.Query("asrama_id")
	list, err := h.service.GetKamarList(c.Context(), tenantID, asramaID)
	if err != nil {
		if err != nil {
			return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data kamar", err)
		}
	}
	return SendSuccess(c, "Data kamar berhasil diambil", list)
}

func (h *akademikHandler) CreateKamar(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields
	var input struct {
		AsramaID  string `json:"asrama_id"`
		Nomor     string `json:"nomor"`
		Kapasitas int    `json:"kapasitas"`
		Status    string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		if err := c.BodyParser(&input); err != nil {
			return SendError(c, fiber.StatusBadRequest, "Invalid request", err)
		}
	}

	// Explicit mapping: DTO → DB Model
	req := model.Kamar{
		TenantID:  tenantID, // From JWT, not user input
		AsramaID:  input.AsramaID,
		Nomor:     input.Nomor,
		Kapasitas: input.Kapasitas,
		Status:    input.Status,
	}

	if err := h.service.CreateKamar(c.Context(), tenantID, &req); err != nil {
		if err := h.service.CreateKamar(c.Context(), tenantID, &req); err != nil {
			return SendError(c, fiber.StatusInternalServerError, "Gagal membuat kamar", err)
		}
	}
	return SendCreated(c, "Kamar created", req)
}

// Penempatan Handlers
func (h *akademikHandler) GetPenempatanList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	list, err := h.service.GetPenempatanList(c.Context(), tenantID)
	if err != nil {
		if err != nil {
			return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data penempatan", err)
		}
	}
	return SendSuccess(c, "Data penempatan berhasil diambil", list)
}

func (h *akademikHandler) CreatePenempatan(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields
	var input struct {
		SantriID     string `json:"santri_id"`
		KamarID      string `json:"kamar_id"`
		TanggalMasuk string `json:"tanggal_masuk"`
		Status       string `json:"status"`
		Keterangan   string `json:"keterangan"`
	}
	if err := c.BodyParser(&input); err != nil {
		if err := c.BodyParser(&input); err != nil {
			return SendError(c, fiber.StatusBadRequest, "Invalid request", err)
		}
	}

	// Explicit mapping: DTO → DB Model
	req := model.Penempatan{
		TenantID:   tenantID, // From JWT, not user input
		SantriID:   input.SantriID,
		KamarID:    input.KamarID,
		Status:     input.Status,
		Keterangan: input.Keterangan,
	}

	// Parse TanggalMasuk if provided
	if input.TanggalMasuk != "" {
		if t, err := time.Parse("2006-01-02", input.TanggalMasuk); err == nil {
			req.TanggalMasuk = t
		}
	}

	if err := h.service.CreatePenempatan(c.Context(), tenantID, &req); err != nil {
		if err := h.service.CreatePenempatan(c.Context(), tenantID, &req); err != nil {
			return SendError(c, fiber.StatusInternalServerError, "Gagal membuat penempatan", err)
		}
	}
	return SendCreated(c, "Penempatan created", req)
}
