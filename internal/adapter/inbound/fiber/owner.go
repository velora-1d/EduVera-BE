package fiber_inbound_adapter

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
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
	ctx := context.Background()

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	if input.Email == "" || input.Password == "" {
		return SendError(c, fiber.StatusBadRequest, "Email dan password wajib diisi.", nil)
	}

	// Find user by email in database
	user, err := h.domain.Auth().GetUserByEmail(ctx, input.Email)
	if err != nil {
		return SendError(c, fiber.StatusUnauthorized, "Email atau password salah. Silakan coba lagi.", err)
	}

	// Check if user is an owner
	if !user.IsOwner {
		return SendError(c, fiber.StatusUnauthorized, "Akun ini bukan akun owner.", nil)
	}

	// Validate password
	if !user.CheckPassword(input.Password) {
		return SendError(c, fiber.StatusUnauthorized, "Email atau password salah. Silakan coba lagi.", nil)
	}

	// Create owner user for JWT (use actual user data)
	ownerUser := &model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     model.RoleSuperAdmin,
		TenantID: "system",
		IsOwner:  true,
	}

	// Generate Token via Auth Domain
	token, expiresAt, err := h.domain.Auth().GenerateToken(ownerUser)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat token. Silakan coba lagi.", err)
	}

	return c.JSON(fiber.Map{
		"access_token": token,
		"expires_at":   expiresAt,
		"user":         ownerUser,
	})
}

// POST /api/v1/owner/impersonate
// Switch view to tenant dashboard for testing/preview
func (h *ownerAdapter) Impersonate(c *fiber.Ctx) error {
	ctx := context.Background()

	var input model.ImpersonateInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	if input.TenantID == "" {
		return SendError(c, fiber.StatusBadRequest, "Tenant ID wajib diisi.", nil)
	}

	// Get tenant info
	tenant, err := h.domain.Tenant().FindByID(ctx, input.TenantID)
	if err != nil {
		return SendError(c, fiber.StatusNotFound, "Tenant tidak ditemukan.", err)
	}

	// Determine view mode based on tenant plan type
	viewMode := input.ViewMode
	if viewMode == "" {
		viewMode = tenant.PlanType // sekolah, pesantren, or hybrid
	}

	// Create impersonation user (owner viewing as tenant admin)
	impersonateUser := &model.User{
		ID:       "owner-impersonate-" + tenant.ID,
		Name:     "Owner (Viewing: " + tenant.Name + ")",
		Email:    "owner@eduvera.id",
		Role:     model.RoleOwner,
		TenantID: tenant.ID,
		IsOwner:  true,
	}

	// Generate Token with impersonation context
	token, expiresAt, err := h.domain.Auth().GenerateToken(impersonateUser)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat token impersonate.", err)
	}

	return c.JSON(fiber.Map{
		"impersonate_token": token,
		"expires_at":        expiresAt,
		"tenant":            tenant,
		"view_mode":         viewMode,
		"is_impersonating":  true,
	})
}

// GET /api/v1/owner/tenants
func (h *ownerAdapter) GetTenants(c *fiber.Ctx) error {
	ctx := context.Background()

	tenants, err := h.domain.Tenant().GetAll(ctx)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat data tenant.", err)
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
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat statistik.", err)
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
		return SendError(c, fiber.StatusNotFound, "Tenant tidak ditemukan.", err)
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
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	// Validate status
	if input.Status != model.TenantStatusActive &&
		input.Status != model.TenantStatusPending &&
		input.Status != model.TenantStatusSuspended {
		return SendError(c, fiber.StatusBadRequest, "Status tidak valid. Pilih: active, pending, atau suspended", nil)
	}

	err := h.domain.Tenant().UpdateStatus(ctx, id, input.Status)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengubah status tenant.", err)
	}

	// Log Admin Action
	_ = h.domain.AuditLog().LogAction(ctx, &model.AuditLogInput{
		AdminID:     "owner-super-admin",
		AdminEmail:  "owner@eduvera.id",
		Action:      model.AuditActionTenantStatusUpdate,
		TargetType:  "tenant",
		TargetID:    id,
		NewValue:    input.Status,
		IPAddress:   c.IP(),
		UserAgent:   string(c.Request().Header.UserAgent()),
		Description: "Tenant status updated to " + input.Status,
	})

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
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat data registrasi.", err)
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
	ctx := context.Background()

	transactions, err := h.domain.SPP().ListAll(ctx)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat data transaksi.", err)
	}

	// Calculate stats
	var totalAmount int64
	var pendingAmount int64
	for _, t := range transactions {
		totalAmount += t.Amount
		if t.Status == "pending" {
			pendingAmount += t.Amount
		}
	}

	return c.JSON(fiber.Map{
		"data": transactions,
		"stats": map[string]interface{}{
			"total_transactions": len(transactions),
			"total_amount":       totalAmount,
			"pending_amount":     pendingAmount,
		},
	})
}

// GET /api/v1/owner/disbursements
func (h *ownerAdapter) GetDisbursements(c *fiber.Ctx) error {
	ctx := context.Background()
	disbursements, err := h.domain.Disbursement().GetAll(ctx)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat data pencairan.", err)
	}

	return c.JSON(fiber.Map{
		"data": disbursements,
	})
}

// POST /api/v1/owner/disbursements/:id/approve
func (h *ownerAdapter) ApproveDisbursement(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	if err := h.domain.Disbursement().Approve(ctx, id); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyetujui pencairan.", err)
	}

	// Log Admin Action
	_ = h.domain.AuditLog().LogAction(ctx, &model.AuditLogInput{
		AdminID:     "owner-super-admin",
		AdminEmail:  "owner@eduvera.id",
		Action:      model.AuditActionDisbursementApprove,
		TargetType:  "disbursement",
		TargetID:    id,
		NewValue:    "approved",
		IPAddress:   c.IP(),
		UserAgent:   string(c.Request().Header.UserAgent()),
		Description: "Disbursement approved",
	})

	return c.JSON(fiber.Map{
		"message": "Disbursement approved successfully",
		"id":      id,
		"status":  "completed",
	})
}

// POST /api/v1/owner/disbursements/:id/reject
func (h *ownerAdapter) RejectDisbursement(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	var input struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	if err := h.domain.Disbursement().Reject(ctx, id, input.Reason); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menolak pencairan.", err)
	}

	// Log Admin Action
	_ = h.domain.AuditLog().LogAction(ctx, &model.AuditLogInput{
		AdminID:     "owner-super-admin",
		AdminEmail:  "owner@eduvera.id",
		Action:      model.AuditActionDisbursementReject,
		TargetType:  "disbursement",
		TargetID:    id,
		NewValue:    "rejected: " + input.Reason,
		IPAddress:   c.IP(),
		UserAgent:   string(c.Request().Header.UserAgent()),
		Description: "Disbursement rejected with reason: " + input.Reason,
	})

	return c.JSON(fiber.Map{
		"message": "Disbursement rejected",
		"id":      id,
		"status":  "rejected",
		"reason":  input.Reason,
	})
}

// GET /api/v1/owner/notifications
func (h *ownerAdapter) GetNotificationLogs(c *fiber.Ctx) error {
	ctx := context.Background()

	notifications, err := h.domain.Notification().GetAll(ctx)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat notifikasi.", err)
	}

	stats, err := h.domain.Notification().GetStats(ctx)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memuat statistik notifikasi.", err)
	}

	return c.JSON(fiber.Map{
		"data":  notifications,
		"stats": stats,
	})
}
