package fiber_inbound_adapter

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"eduvera/internal/domain"
	"eduvera/internal/model"
	inbound_port "eduvera/internal/port/inbound"
)

type onboardingAdapter struct {
	domain domain.Domain
}

func NewOnboardingAdapter(domain domain.Domain) inbound_port.OnboardingHttpPort {
	return &onboardingAdapter{
		domain: domain,
	}
}

// RegisterInput is the request body for the register endpoint
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	WhatsApp string `json:"whatsapp"`
	Password string `json:"password"`
}

// InstitutionInput is the request body for the institution endpoint
type InstitutionInput struct {
	TenantID        string `json:"tenant_id"`
	InstitutionName string `json:"institution_name"`
	InstitutionType string `json:"institution_type"`
	PlanType        string `json:"plan_type"`
	Address         string `json:"address"`
}

// SubdomainInput is the request body for the subdomain endpoint
type SubdomainInput struct {
	TenantID  string `json:"tenant_id"`
	Subdomain string `json:"subdomain"`
}

// BankAccountInput is the request body for the bank account endpoint
type BankAccountInput struct {
	TenantID      string `json:"tenant_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountHolder string `json:"account_holder"`
}

// ConfirmInput is the request body for the confirm endpoint
type ConfirmInput struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

// POST /api/v1/onboarding/register
func (h *onboardingAdapter) Register(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required",
		})
	}

	if len(input.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters",
		})
	}

	// Register user (without tenant for now - tenant will be created in institution step)
	user, err := h.domain.Auth().Register(ctx, &model.UserInput{
		Name:     input.Name,
		Email:    input.Email,
		WhatsApp: input.WhatsApp,
		Password: input.Password,
		Role:     model.RoleAdmin,
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"message":   "Admin account created",
		"user_id":   user.ID,
		"next_step": "/api/v1/onboarding/institution",
	})
}

// POST /api/v1/onboarding/institution
func (h *onboardingAdapter) Institution(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input InstitutionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create tenant with institution info
	tenant, err := h.domain.Tenant().Create(ctx, &model.TenantInput{
		Name:            input.InstitutionName,
		Subdomain:       generateSubdomain(input.InstitutionName),
		PlanType:        input.PlanType,
		InstitutionType: input.InstitutionType,
		Address:         input.Address,
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":              "success",
		"message":             "Institution created",
		"tenant_id":           tenant.ID,
		"suggested_subdomain": tenant.Subdomain,
		"next_step":           "/api/v1/onboarding/subdomain",
	})
}

// POST /api/v1/onboarding/subdomain
func (h *onboardingAdapter) Subdomain(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input SubdomainInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check subdomain availability
	exists, err := h.domain.Tenant().SubdomainExists(ctx, input.Subdomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check subdomain",
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":     "Subdomain already taken",
			"available": false,
		})
	}

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Subdomain available",
		"subdomain": input.Subdomain,
		"available": true,
		"next_step": "/api/v1/onboarding/bank-account",
	})
}

// POST /api/v1/onboarding/bank-account
func (h *onboardingAdapter) BankAccount(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input BankAccountInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if input.BankName == "" || input.AccountNumber == "" || input.AccountHolder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bank name, account number, and account holder are required",
		})
	}

	// Update tenant with bank account
	err := h.domain.Tenant().UpdateBankAccount(ctx,
		input.TenantID,
		input.BankName,
		input.AccountNumber,
		input.AccountHolder,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save bank account",
		})
	}

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Bank account saved",
		"next_step": "/api/v1/onboarding/confirm",
	})
}

// POST /api/v1/onboarding/confirm
func (h *onboardingAdapter) Confirm(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input ConfirmInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Activate tenant
	err := h.domain.Tenant().Activate(ctx, input.TenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate account",
		})
	}

	// Get tenant for subdomain info
	tenant, err := h.domain.Tenant().FindByID(ctx, input.TenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tenant info",
		})
	}

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Registration complete! Your account is now active.",
		"subdomain": tenant.Subdomain + ".eduvera.id",
		"login_url": "https://" + tenant.Subdomain + ".eduvera.id/login",
	})
}

// GET /api/v1/onboarding/status/:id
func (h *onboardingAdapter) Status(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	tenantID := c.Params("id")
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID required",
		})
	}

	tenant, err := h.domain.Tenant().FindByID(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	return c.JSON(fiber.Map{
		"tenant_id": tenant.ID,
		"name":      tenant.Name,
		"subdomain": tenant.Subdomain,
		"plan_type": tenant.PlanType,
		"status":    tenant.Status,
		"is_active": tenant.Status == model.TenantStatusActive,
	})
}

// Helper function to generate subdomain from institution name
func generateSubdomain(name string) string {
	subdomain := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			subdomain += string(c)
		} else if c >= 'A' && c <= 'Z' {
			subdomain += string(c + 32) // lowercase
		}
	}
	if len(subdomain) > 30 {
		subdomain = subdomain[:30]
	}
	return subdomain
}
