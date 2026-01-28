package fiber_inbound_adapter

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"

	"eduvera/internal/domain"
	"eduvera/internal/model"
	inbound_port "eduvera/internal/port/inbound"
)

type ownerAdapter struct {
	domain domain.Domain
}

func NewOwnerAdapter(domain domain.Domain) inbound_port.OwnerHttpPort {
	return &ownerAdapter{
		domain: domain,
	}
}

// POST /api/v1/owner/login
func (h *ownerAdapter) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate against .env credentials
	envEmail := os.Getenv("OWNER_EMAIL")
	envPassword := os.Getenv("OWNER_PASSWORD")

	if envEmail == "" || envPassword == "" {
		// Fallback for safety if env not set
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Owner credentials not configured",
		})
	}

	if input.Email != envEmail || input.Password != envPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Create mock owner user for JWT
	ownerUser := &model.User{
		ID:       "owner-super-admin",
		Name:     "Owner",
		Email:    "owner@eduvera.id",
		Role:     model.RoleSuperAdmin,
		TenantID: "system",
	}

	// Generate Token via Auth Domain
	token, expiresAt, err := h.domain.Auth().GenerateToken(ownerUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token":      token,
		"expires_at": expiresAt,
		"user":       ownerUser,
	})
}

// GET /api/v1/owner/tenants
func (h *ownerAdapter) GetTenants(c *fiber.Ctx) error {
	ctx := context.Background()

	tenants, err := h.domain.Tenant().GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": tenants,
	})
}

// GET /api/v1/owner/stats
func (h *ownerAdapter) GetStats(c *fiber.Ctx) error {
	ctx := context.Background()

	tenants, err := h.domain.Tenant().GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalTenants := len(tenants)
	activeTenants := 0
	totalRevenue := int64(0)

	for _, t := range tenants {
		if t.Status == model.TenantStatusActive {
			activeTenants++
		}
		// Estimate revenue (simplified)
		if t.Status == model.TenantStatusActive {
			// Mock calculation based on plan type
			switch t.InstitutionType {
			case "sekolah":
				totalRevenue += 500000
			case "pesantren":
				totalRevenue += 350000
			case "hybrid":
				totalRevenue += 750000
			}
		}
	}

	return c.JSON(fiber.Map{
		"total_tenants":  totalTenants,
		"active_tenants": activeTenants,
		"total_revenue":  totalRevenue,
	})
}

// GET /api/v1/owner/tenants/:id
func (h *ownerAdapter) GetTenantDetail(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	tenant, err := h.domain.Tenant().FindByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": tenant,
	})
}

// PUT /api/v1/owner/tenants/:id/status
func (h *ownerAdapter) UpdateTenantStatus(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	if input.Status != model.TenantStatusActive &&
		input.Status != model.TenantStatusPending &&
		input.Status != model.TenantStatusSuspended {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be: active, pending, or suspended",
		})
	}

	err := h.domain.Tenant().UpdateStatus(ctx, id, input.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status updated successfully",
		"status":  input.Status,
	})
}

// GET /api/v1/owner/registrations
func (h *ownerAdapter) GetRegistrations(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get all tenants sorted by created_at desc (registration logs)
	tenants, err := h.domain.Tenant().GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Map to registration format
	registrations := make([]map[string]interface{}, 0)
	for _, t := range tenants {
		registrations = append(registrations, map[string]interface{}{
			"id":               t.ID,
			"name":             t.Name,
			"subdomain":        t.Subdomain,
			"plan_type":        t.PlanType,
			"institution_type": t.InstitutionType,
			"status":           t.Status,
			"registered_at":    t.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"data": registrations,
	})
}

// GET /api/v1/owner/transactions
func (h *ownerAdapter) GetSPPTransactions(c *fiber.Ctx) error {
	// For now, return mock data since SPP transactions table might not be fully implemented
	// This will be replaced with actual database query when payment gateway is integrated

	transactions := []map[string]interface{}{
		{
			"id":          "TRX-001",
			"tenant_name": "SD Negeri 1 Jakarta",
			"student":     "Ahmad Rizki",
			"amount":      150000,
			"type":        "SPP Juli 2024",
			"status":      "success",
			"created_at":  "2024-07-15T10:00:00Z",
		},
		{
			"id":          "TRX-002",
			"tenant_name": "SMP Al-Azhar",
			"student":     "Siti Nurhaliza",
			"amount":      200000,
			"type":        "SPP Juli 2024",
			"status":      "success",
			"created_at":  "2024-07-14T14:30:00Z",
		},
		{
			"id":          "TRX-003",
			"tenant_name": "SMK Teknologi",
			"student":     "Budi Santoso",
			"amount":      175000,
			"type":        "SPP Juli 2024",
			"status":      "pending",
			"created_at":  "2024-07-13T09:15:00Z",
		},
	}

	return c.JSON(fiber.Map{
		"data": transactions,
		"stats": map[string]interface{}{
			"total_transactions": 156,
			"total_amount":       23400000,
			"pending_amount":     1750000,
		},
	})
}

// GET /api/v1/owner/disbursements
func (h *ownerAdapter) GetDisbursements(c *fiber.Ctx) error {
	// Mock data for disbursements
	disbursements := []map[string]interface{}{
		{
			"id":             "DSB-001",
			"tenant_name":    "SD Negeri 1 Jakarta",
			"account_name":   "Yayasan Pendidikan Mandiri",
			"bank":           "BCA",
			"account_number": "1234567890",
			"amount":         5000000,
			"status":         "pending",
			"requested_at":   "2024-07-20T10:00:00Z",
		},
		{
			"id":             "DSB-002",
			"tenant_name":    "SMP Al-Azhar",
			"account_name":   "Yayasan Al-Azhar",
			"bank":           "Mandiri",
			"account_number": "0987654321",
			"amount":         7500000,
			"status":         "approved",
			"requested_at":   "2024-07-19T14:30:00Z",
		},
	}

	return c.JSON(fiber.Map{
		"data": disbursements,
		"stats": map[string]interface{}{
			"pending_count":   3,
			"pending_amount":  15000000,
			"approved_count":  12,
			"approved_amount": 120000000,
		},
	})
}

// POST /api/v1/owner/disbursements/:id/approve
func (h *ownerAdapter) ApproveDisbursement(c *fiber.Ctx) error {
	id := c.Params("id")

	// TODO: Update disbursement status in database
	// For now, return success response

	return c.JSON(fiber.Map{
		"message": "Disbursement approved successfully",
		"id":      id,
		"status":  "approved",
	})
}

// POST /api/v1/owner/disbursements/:id/reject
func (h *ownerAdapter) RejectDisbursement(c *fiber.Ctx) error {
	id := c.Params("id")

	var input struct {
		Reason string `json:"reason"`
	}
	c.BodyParser(&input)

	// TODO: Update disbursement status in database
	// For now, return success response

	return c.JSON(fiber.Map{
		"message": "Disbursement rejected",
		"id":      id,
		"status":  "rejected",
		"reason":  input.Reason,
	})
}

// GET /api/v1/owner/notifications
func (h *ownerAdapter) GetNotificationLogs(c *fiber.Ctx) error {
	// Mock data for notification logs
	notifications := []map[string]interface{}{
		{
			"id":         "NOT-001",
			"type":       "telegram",
			"recipient":  "Owner Channel",
			"message":    "Pendaftaran baru: SD Negeri 1 Jakarta",
			"status":     "sent",
			"created_at": "2024-07-20T10:00:00Z",
		},
		{
			"id":         "NOT-002",
			"type":       "email",
			"recipient":  "admin@sdnegeri1.sch.id",
			"message":    "Selamat datang di EduVera!",
			"status":     "sent",
			"created_at": "2024-07-20T10:01:00Z",
		},
		{
			"id":         "NOT-003",
			"type":       "whatsapp",
			"recipient":  "+6281234567890",
			"message":    "Pembayaran SPP berhasil - Ahmad Rizki",
			"status":     "failed",
			"created_at": "2024-07-19T14:30:00Z",
		},
	}

	return c.JSON(fiber.Map{
		"data": notifications,
		"stats": map[string]interface{}{
			"total_sent":   256,
			"total_failed": 12,
		},
	})
}
