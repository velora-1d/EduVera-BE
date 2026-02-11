package fiber_inbound_adapter

import (
	"context"
	"net/url"
	"strings"

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

	// SECURITY: Get tenant_id from JWT context, not from query params
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid. Silakan login kembali.", nil)
	}

	transactions, err := h.domain.SPP().ListByTenant(ctx, tenantID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat data SPP.", err)
	}

	return SendSuccess(c, "Data SPP berhasil dimuat", fiber.Map{
		"data": transactions,
	})
}

// POST /api/v1/tenant/spp
func (h *sppAdapter) Create(c *fiber.Ctx) error {
	ctx := context.Background()

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid. Silakan login kembali.", nil)
	}

	var input struct {
		StudentID   string `json:"student_id"`
		StudentName string `json:"student_name"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	spp := &model.SPPTransaction{
		TenantID:    tenantID, // From JWT, not user input
		StudentID:   input.StudentID,
		StudentName: input.StudentName,
		Amount:      input.Amount,
		Description: input.Description,
	}

	if err := h.domain.SPP().Create(ctx, spp); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat tagihan SPP.", err)
	}

	return SendSuccess(c, "Tagihan berhasil dibuat", spp)
}

// POST /api/v1/tenant/spp/:id/pay
func (h *sppAdapter) RecordPayment(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid. Silakan coba lagi.",
		})
	}

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses tidak valid",
		})
	}

	if err := h.domain.SPP().RecordPayment(ctx, tenantID, id, input.PaymentMethod); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mencatat pembayaran.", err)
	}

	return SendSuccess(c, "Pembayaran berhasil dicatat", fiber.Map{"status": "paid"})
}

// GET /api/v1/tenant/spp/stats
func (h *sppAdapter) GetStats(c *fiber.Ctx) error {
	ctx := context.Background()

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid. Silakan login kembali.", nil)
	}

	stats, err := h.domain.SPP().GetStats(ctx, tenantID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat statistik.", err)
	}

	return SendSuccess(c, "Statistik SPP berhasil dimuat", stats)
}

// PUT /api/v1/tenant/spp/:id
func (h *sppAdapter) Update(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		StudentName string `json:"student_name"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"` // Format: 2024-01-31
		Period      string `json:"period"`   // Format: 2024-01
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid. Silakan coba lagi.",
		})
	}

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid", nil)
	}

	if err := h.domain.SPP().Update(ctx, tenantID, id, input.StudentName, input.Amount, input.Description, input.DueDate, input.Period); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memperbarui tagihan.", err)
	}

	return SendSuccess(c, "Tagihan berhasil diperbarui", nil)
}

// DELETE /api/v1/tenant/spp/:id
func (h *sppAdapter) Delete(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid", nil)
	}

	if err := h.domain.SPP().Delete(ctx, tenantID, id); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menghapus tagihan.", err)
	}

	return SendSuccess(c, "Tagihan berhasil dihapus", nil)
}

// POST /api/v1/tenant/spp/:id/upload-proof
func (h *sppAdapter) UploadProof(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		ProofURL string `json:"proof_url"` // URL dari uploaded image
	}

	if err := c.BodyParser(&input); err != nil || input.ProofURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL bukti pembayaran wajib diisi.",
		})
	}

	// SECURITY: Validate URL format and only allow safe URLs
	if !isValidProofURL(input.ProofURL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL tidak valid. Hanya URL HTTPS yang diizinkan.",
		})
	}

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid", nil)
	}

	if err := h.domain.SPP().UploadProof(ctx, tenantID, id, input.ProofURL); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan bukti pembayaran.", err)
	}

	return SendSuccess(c, "Bukti pembayaran berhasil diupload", nil)
}

// isValidProofURL validates that the URL is safe to store
func isValidProofURL(proofURL string) bool {
	// Parse URL
	parsed, err := url.Parse(proofURL)
	if err != nil {
		return false
	}

	// Only allow HTTPS
	if parsed.Scheme != "https" {
		return false
	}

	// Block dangerous protocols and paths
	if strings.Contains(proofURL, "javascript:") ||
		strings.Contains(proofURL, "data:") ||
		strings.Contains(proofURL, "..") {
		return false
	}

	// Allowed domains for uploaded images (adjust based on your storage)
	allowedDomains := []string{
		"storage.googleapis.com",
		"res.cloudinary.com",
		"s3.amazonaws.com",
		"eduvera.ve-lora.my.id",
	}

	for _, domain := range allowedDomains {
		if strings.HasSuffix(parsed.Host, domain) {
			return true
		}
	}

	// If no allowed domain matches, reject
	return false
}

// POST /api/v1/tenant/spp/:id/confirm
func (h *sppAdapter) ConfirmPayment(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	// Get user ID from JWT context
	userID := c.Locals("user_id")
	confirmedBy := ""
	if userID != nil {
		confirmedBy = userID.(string)
	}

	var input struct {
		PaymentMethod string `json:"payment_method"` // cash, transfer
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid. Silakan coba lagi.",
		})
	}

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid", nil)
	}

	if err := h.domain.SPP().ConfirmPayment(ctx, tenantID, id, confirmedBy, input.PaymentMethod); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal konfirmasi pembayaran.", err)
	}

	return SendSuccess(c, "Pembayaran berhasil dikonfirmasi", fiber.Map{"status": "paid"})
}

// GET /api/v1/tenant/spp/overdue
func (h *sppAdapter) ListOverdue(c *fiber.Ctx) error {
	ctx := context.Background()

	// SECURITY: Get tenant_id from JWT context
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return SendError(c, fiber.StatusUnauthorized, "Akses tidak valid. Silakan login kembali.", nil)
	}

	transactions, err := h.domain.SPP().ListOverdue(ctx, tenantID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat data tunggakan.", err)
	}

	return SendSuccess(c, "Data tunggakan berhasil dimuat", transactions)
}

// POST /api/v1/owner/invoices/generate
func (h *sppAdapter) GenerateManual(c *fiber.Ctx) error {
	ctx := context.Background()

	// Use goroutine to avoid timeout as it has delays
	go func() {
		_ = h.domain.SPP().GenerateInvoices(ctx)
	}()

	return c.JSON(fiber.Map{
		"message": "Proses pembuatan invoice dimulai di background (Anti-Spam active)",
	})
}

// POST /api/v1/owner/invoices/broadcast
func (h *sppAdapter) BroadcastOverdueManual(c *fiber.Ctx) error {
	ctx := context.Background()

	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses tidak valid",
		})
	}

	// Use goroutine to avoid timeout
	go func() {
		_ = h.domain.SPP().BroadcastOverdue(ctx, tenantID)
	}()

	return c.JSON(fiber.Map{
		"message": "Proses broadcast tunggakan dimulai di background (Anti-Spam active)",
	})
}
