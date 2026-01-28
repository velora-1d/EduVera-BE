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
