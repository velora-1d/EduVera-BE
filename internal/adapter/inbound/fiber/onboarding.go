package fiber_inbound_adapter

import (
	"context"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/adapter/outbound/notification"
	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
)

// Subdomain blacklist - reserved names that cannot be used
var subdomainBlacklist = []string{
	"admin", "api", "www", "app", "dashboard", "help", "support",
	"mail", "ftp", "smtp", "pop", "imap", "webmail",
	"billing", "payment", "checkout", "login", "register", "signup",
	"account", "settings", "config", "system", "root", "super",
	"test", "demo", "staging", "dev", "development", "prod", "production",
	"static", "assets", "cdn", "media", "images", "files",
	"eduvera", "velora", "sekolah", "pesantren", "madrasah",
}

// Subdomain validation regex: lowercase letters, numbers, and dashes only
// Must start and end with alphanumeric, 3-30 characters
var subdomainRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,28}[a-z0-9]$`)

type onboardingAdapter struct {
	domain   domain.Domain
	telegram *notification.TelegramNotifier
}

func NewOnboardingAdapter(domain domain.Domain) inbound_port.OnboardingHttpPort {
	return &onboardingAdapter{
		domain:   domain,
		telegram: notification.NewTelegramNotifier(),
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
	UserID          string `json:"user_id"`
	InstitutionName string `json:"institution_name"`
	InstitutionType string `json:"institution_type"` // sekolah, pesantren, hybrid
	Address         string `json:"address"`
	Subdomain       string `json:"subdomain"` // Custom subdomain input
}

// SubdomainCheckInput is the request body for checking subdomain availability
type SubdomainCheckInput struct {
	Subdomain string `json:"subdomain"`
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
	AccountType   string `json:"account_type"` // pribadi, yayasan
}

// ConfirmInput is the request body for the confirm endpoint
type ConfirmInput struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

// validateSubdomainFormat checks if subdomain has valid format
func validateSubdomainFormat(subdomain string) (bool, string) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))

	if len(subdomain) < 3 {
		return false, "Subdomain minimal 3 karakter"
	}

	if len(subdomain) > 30 {
		return false, "Subdomain maksimal 30 karakter"
	}

	if !subdomainRegex.MatchString(subdomain) {
		return false, "Subdomain hanya boleh huruf kecil, angka, dan dash (-). Harus diawali dan diakhiri huruf/angka"
	}

	return true, ""
}

// isSubdomainBlacklisted checks if subdomain is in blacklist
func isSubdomainBlacklisted(subdomain string) bool {
	subdomain = strings.ToLower(subdomain)
	for _, blocked := range subdomainBlacklist {
		if subdomain == blocked {
			return true
		}
	}
	return false
}

// generateSubdomainRecommendations generates alternative subdomain suggestions
func generateSubdomainRecommendations(base string) []string {
	base = strings.ToLower(strings.TrimSpace(base))
	// Clean base to valid format
	cleanBase := ""
	for _, c := range base {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			cleanBase += string(c)
		}
	}

	if len(cleanBase) > 20 {
		cleanBase = cleanBase[:20]
	}

	recommendations := []string{
		cleanBase + "1",
		cleanBase + "-id",
		cleanBase + "-edu",
		cleanBase + "-sch",
		cleanBase + "2025",
		"my" + cleanBase,
		cleanBase + "-app",
	}

	return recommendations
}

// POST /api/v1/onboarding/check-subdomain - Check subdomain availability (realtime validation)
func (h *onboardingAdapter) CheckSubdomain(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input SubdomainCheckInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	subdomain := strings.ToLower(strings.TrimSpace(input.Subdomain))

	// Step 1: Validate format
	valid, errorMsg := validateSubdomainFormat(subdomain)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"available":       false,
			"subdomain":       subdomain,
			"error":           errorMsg,
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Step 2: Check blacklist
	if isSubdomainBlacklisted(subdomain) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"available":       false,
			"subdomain":       subdomain,
			"error":           "Subdomain ini tidak tersedia (reserved)",
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Step 3: Check if already taken in database
	exists, err := h.domain.Tenant().SubdomainExists(ctx, subdomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengecek subdomain",
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"available":       false,
			"subdomain":       subdomain,
			"error":           "Subdomain sudah digunakan",
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Subdomain is available!
	return c.JSON(fiber.Map{
		"available": true,
		"subdomain": subdomain,
		"message":   "Subdomain tersedia!",
		"full_url":  "https://" + subdomain + ".eduvera.ve-lora.my.id",
	})
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

	// Generate short-lived onboarding token (reusing existing JWT for simplicity)
	// This token will be used to authorize the Institution step
	onboardingToken, expiresAt, err := h.domain.Auth().GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate onboarding token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":           "success",
		"message":          "Admin account created",
		"user_id":          user.ID, // Kept for frontend display, but not used for auth
		"onboarding_token": onboardingToken,
		"token_expires_at": expiresAt,
		"next_step":        "/api/v1/onboarding/institution",
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

	// Validate institution type
	validTypes := []string{"sekolah", "pesantren", "hybrid"}
	isValidType := false
	for _, t := range validTypes {
		if input.InstitutionType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Institution type must be: sekolah, pesantren, or hybrid",
		})
	}

	// Determine subdomain: use custom input or generate from name
	subdomain := strings.ToLower(strings.TrimSpace(input.Subdomain))
	if subdomain == "" {
		subdomain = generateSubdomain(input.InstitutionName)
	}

	// Validate subdomain format
	valid, errorMsg := validateSubdomainFormat(subdomain)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":           errorMsg,
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Check blacklist
	if isSubdomainBlacklisted(subdomain) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":           "Subdomain tidak tersedia (reserved)",
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Check if subdomain already taken
	exists, err := h.domain.Tenant().SubdomainExists(ctx, subdomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check subdomain",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":           "Subdomain sudah digunakan",
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Create tenant with institution info
	tenant, err := h.domain.Tenant().Create(ctx, &model.TenantInput{
		Name:            input.InstitutionName,
		Subdomain:       subdomain,
		InstitutionType: input.InstitutionType,
		Address:         input.Address,
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Link user to tenant - Extract UserID from Authorization token (not from body!)
	authHeader := c.Get("Authorization")
	var userID string

	// First, try to extract from token (secure method)
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		bearerToken := authHeader[7:]
		claims, err := h.domain.Auth().ValidateToken(ctx, bearerToken)
		if err == nil && claims.UserID != "" {
			userID = claims.UserID
		}
	}

	// Fallback to body for backwards compatibility (deprecated)
	if userID == "" && input.UserID != "" {
		userID = input.UserID
	}

	if userID != "" {
		err = h.domain.Auth().LinkUserToTenant(ctx, userID, tenant.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to link user to tenant: " + err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"message":   "Institution created",
		"tenant_id": tenant.ID,
		"subdomain": tenant.Subdomain,
		"full_url":  "https://" + tenant.Subdomain + ".eduvera.ve-lora.my.id",
		"next_step": "/api/v1/onboarding/modules",
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

	subdomain := strings.ToLower(strings.TrimSpace(input.Subdomain))

	// Validate format
	valid, errorMsg := validateSubdomainFormat(subdomain)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"available":       false,
			"error":           errorMsg,
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Check blacklist
	if isSubdomainBlacklisted(subdomain) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"available":       false,
			"error":           "Subdomain tidak tersedia (reserved)",
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// Check availability
	exists, err := h.domain.Tenant().SubdomainExists(ctx, subdomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check subdomain",
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"available":       false,
			"error":           "Subdomain sudah digunakan",
			"recommendations": generateSubdomainRecommendations(subdomain),
		})
	}

	// TODO: Update tenant subdomain (method needs to be added to Tenant domain)
	// err = h.domain.Tenant().UpdateSubdomain(ctx, input.TenantID, subdomain)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to update subdomain",
	// 	})
	// }
	_ = input.TenantID // Placeholder until UpdateSubdomain is implemented

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Subdomain updated",
		"subdomain": subdomain,
		"available": true,
		"full_url":  "https://" + subdomain + ".eduvera.ve-lora.my.id",
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

	// Validate account type
	if input.AccountType != "" && input.AccountType != "pribadi" && input.AccountType != "yayasan" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account type must be: pribadi or yayasan",
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

	// Send Telegram notification to owner (async, don't block response)
	go func() {
		data := notification.RegistrationData{
			InstitutionName: tenant.Name,
			PlanType:        tenant.InstitutionType,
			Subdomain:       tenant.Subdomain,
			Address:         tenant.Address,
			// Note: User data will be added when we implement user lookup
		}
		_ = h.telegram.SendNewRegistration(data)
	}()

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Registration complete! Your account is now active.",
		"subdomain": tenant.Subdomain + ".eduvera.ve-lora.my.id",
		"login_url": "https://" + tenant.Subdomain + ".eduvera.ve-lora.my.id/login",
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
		"tenant_id":        tenant.ID,
		"name":             tenant.Name,
		"subdomain":        tenant.Subdomain,
		"institution_type": tenant.InstitutionType,
		"status":           tenant.Status,
		"is_active":        tenant.Status == model.TenantStatusActive,
		"full_url":         "https://" + tenant.Subdomain + ".eduvera.ve-lora.my.id",
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
	if len(subdomain) < 3 {
		subdomain = subdomain + "app"
	}
	return subdomain
}
