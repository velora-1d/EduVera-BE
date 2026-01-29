package fiber_inbound_adapter

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	"prabogo/internal/model"
)

type sppAdapter struct {
	domain domain.Domain
}

func NewSPPAdapter(domain domain.Domain) *sppAdapter {
	return &sppAdapter{
		domain: domain,
	}
}

// GET /api/v1/tenant/spp
func (h *sppAdapter) List(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tenant wajib diisi.",
		})
	}

	transactions, err := h.domain.SPP().ListByTenant(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memuat data SPP. " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": transactions,
	})
}

// POST /api/v1/tenant/spp
func (h *sppAdapter) Create(c *fiber.Ctx) error {
	ctx := context.Background()

	var input struct {
		TenantID    string `json:"tenant_id"`
		StudentID   string `json:"student_id"`
		StudentName string `json:"student_name"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid. Silakan coba lagi.",
		})
	}

	spp := &model.SPPTransaction{
		TenantID:    input.TenantID,
		StudentID:   input.StudentID,
		StudentName: input.StudentName,
		Amount:      input.Amount,
		Description: input.Description,
	}

	if err := h.domain.SPP().Create(ctx, spp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat tagihan SPP. " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tagihan berhasil dibuat",
		"data":    spp,
	})
}

// POST /api/v1/tenant/spp/:id/pay
func (h *sppAdapter) RecordPayment(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		PaymentMethod string `json:"payment_method"`
	}
	c.BodyParser(&input)

	if err := h.domain.SPP().RecordPayment(ctx, id, input.PaymentMethod); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mencatat pembayaran. " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pembayaran berhasil dicatat",
		"status":  "paid",
	})
}

// GET /api/v1/tenant/spp/stats
func (h *sppAdapter) GetStats(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "tenant_id is required",
		})
	}

	stats, err := h.domain.SPP().GetStats(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": stats,
	})
}
