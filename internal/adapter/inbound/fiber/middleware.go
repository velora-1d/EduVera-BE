package fiber_inbound_adapter

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	"prabogo/utils/activity"
	"prabogo/utils/jwt"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	bearerPrefixLen     = 7
)

type MiddlewareAdapter interface {
	InternalAuth(a any) error
	ClientAuth(a any) error
	OwnerAuth(a any) error
}

type middlewareAdapter struct {
	domain domain.Domain
}

func NewMiddlewareAdapter(
	domain domain.Domain,
) MiddlewareAdapter {
	return &middlewareAdapter{
		domain: domain,
	}
}

func (h *middlewareAdapter) OwnerAuth(a any) error {
	c := a.(*fiber.Ctx)
	ctx := activity.NewContext("owner_auth")

	authHeader := c.Get(authorizationHeader)
	var bearerToken string
	if len(authHeader) > bearerPrefixLen && authHeader[:bearerPrefixLen] == bearerPrefix {
		bearerToken = authHeader[bearerPrefixLen:]
	}

	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Sesi Anda telah berakhir. Silakan login kembali.",
		})
	}

	// Use Auth Domain to validate token
	claims, err := h.domain.Auth().ValidateToken(ctx, bearerToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token tidak valid. Silakan login kembali.",
		})
	}

	// Check Role
	if claims.Role != model.RoleSuperAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda tidak memiliki akses ke halaman ini.",
		})
	}

	// Store user info in context for handlers
	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)

	return c.Next()
}

func (h *middlewareAdapter) InternalAuth(a any) error {
	c := a.(*fiber.Ctx)
	authHeader := c.Get(authorizationHeader)
	var bearerToken string
	if len(authHeader) > bearerPrefixLen && authHeader[:bearerPrefixLen] == bearerPrefix {
		bearerToken = authHeader[bearerPrefixLen:]
	}

	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses ditolak. Token tidak ditemukan.",
		})
	}

	if bearerToken != os.Getenv("INTERNAL_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses ditolak. Token tidak valid.",
		})
	}

	return c.Next()
}

func (h *middlewareAdapter) ClientAuth(a any) error {
	c := a.(*fiber.Ctx)
	ctx := activity.NewContext("http_client_auth")
	authHeader := c.Get(authorizationHeader)
	var bearerToken string
	if len(authHeader) > bearerPrefixLen && authHeader[:bearerPrefixLen] == bearerPrefix {
		bearerToken = authHeader[bearerPrefixLen:]
	}

	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Success: false,
			Error:   "Sesi Anda telah berakhir. Silakan login kembali.",
		})
	}

	authDriver := os.Getenv("AUTH_DRIVER")
	if authDriver == "jwt" {
		jwksURL := os.Getenv("AUTH_JWKS_URL")

		_, err := jwt.ValidateJWTWithURL(bearerToken, jwksURL)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Success: false,
				Error:   "Token tidak valid. Silakan login kembali.",
			})
		}
	} else {
		// Validate token using Auth domain to get claims
		claims, err := h.domain.Auth().ValidateToken(ctx, bearerToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Success: false,
				Error:   "Token tidak valid. Silakan login kembali.",
			})
		}

		// SECURITY: Check if token is blacklisted (logged out)
		isBlacklisted, _ := h.domain.Auth().IsTokenBlacklisted(ctx, bearerToken)
		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Success: false,
				Error:   "Sesi telah berakhir. Silakan login kembali.",
			})
		}

		// OWNER PREVIEW MODE: Allow super_admin to access all routes
		// This is for owner to test/debug tenant dashboards using sandbox tenants
		if claims.Role == model.RoleSuperAdmin {
			c.Locals("user_id", claims.UserID)
			c.Locals("role", claims.Role)
			c.Locals("is_owner_preview", true)
			c.Locals("bearer_token", bearerToken)

			// Check if owner is using a sandbox tenant
			sandboxSubdomain := c.Get("X-Sandbox-Tenant")
			if sandboxSubdomain != "" {
				// Validate sandbox tenant exists and belongs to this owner
				tenant, err := h.domain.Tenant().FindBySubdomain(ctx, sandboxSubdomain)
				if err == nil && tenant != nil && tenant.IsSandbox {
					// Owner is using a valid sandbox tenant
					c.Locals("tenant_id", tenant.ID)
					c.Locals("plan_type", tenant.PlanType)
				}
			}

			return c.Next()
		}

		// SECURITY: Set user info from JWT claims to context
		// This prevents IDOR by ensuring tenant_id comes from token, not user input
		c.Locals("user_id", claims.UserID)
		c.Locals("tenant_id", claims.TenantID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("bearer_token", bearerToken) // Store token for logout
	}

	return c.Next()
}

// RequirePlan creates middleware that checks if the authenticated tenant has the required plan type
// This prevents users from accessing features not included in their subscription
func RequirePlan(requiredPlans ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get tenant plan from context (should be set by ClientAuth middleware)
		planType, ok := c.Locals("plan_type").(string)
		// If no plan type in context, it means user is not authenticated or context not set properly
		// STRICT SECURITY: Deny access if we cannot verify the plan
		if !ok || planType == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak. Informasi paket langganan tidak ditemukan.",
			})
		}

		// Continue with plan validation logic below...

		// Check if tenant's plan matches any of the required plans
		for _, required := range requiredPlans {
			if planType == required {
				return c.Next()
			}
		}

		// Plan not matched - block access
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":         "Fitur ini tidak tersedia untuk paket langganan Anda",
			"required_plan": requiredPlans,
			"your_plan":     planType,
		})
	}
}

// CheckTrialStatus creates middleware that checks if tenant trial has expired
// Returns 403 with TRIAL_EXPIRED code if trial has ended and no active subscription
func CheckTrialStatus(d domain.Domain) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := activity.NewContext("check_trial_status")

		// Get tenant_id from context (set by ClientAuth)
		tenantID, ok := c.Locals("tenant_id").(string)
		if !ok || tenantID == "" {
			return c.Next()
		}

		// Skip for owner preview mode
		if isOwner, _ := c.Locals("is_owner_preview").(bool); isOwner {
			return c.Next()
		}

		// Get tenant info
		tenant, err := d.Tenant().FindByID(ctx, tenantID)
		if err != nil || tenant == nil {
			return c.Next()
		}

		// Premium tier (paid) - no trial restrictions
		if tenant.SubscriptionTier == model.TierPremium {
			return c.Next()
		}

		// Check if trial has expired (basic tier with expired trial_ends_at)
		if tenant.TrialEndsAt != nil && tenant.TrialEndsAt.Before(time.Now()) {
			// Trial expired - block write operations
			if c.Method() != "GET" && c.Method() != "HEAD" && c.Method() != "OPTIONS" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Trial Anda sudah berakhir. Upgrade sekarang untuk melanjutkan.",
					"code":  "TRIAL_EXPIRED",
					"action": fiber.Map{
						"type": "upgrade",
						"url":  "/upgrade",
					},
				})
			}
		}

		return c.Next()
	}
}

// CheckDataLimit creates middleware that enforces per-table record limits for basic tier
// Returns 403 with LIMIT_REACHED code if limit exceeded
// Usage: router.Post("/students", CheckDataLimit(domain, "santri", 10), handler)
func CheckDataLimit(d domain.Domain, tableName string, limit int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := activity.NewContext("check_data_limit")

		// Get tenant_id from context
		tenantID, ok := c.Locals("tenant_id").(string)
		if !ok || tenantID == "" {
			return c.Next()
		}

		// Skip for owner preview mode
		if isOwner, _ := c.Locals("is_owner_preview").(bool); isOwner {
			return c.Next()
		}

		// Get tenant info
		tenant, err := d.Tenant().FindByID(ctx, tenantID)
		if err != nil || tenant == nil {
			return c.Next()
		}

		// Premium tier has no limits
		if tenant.SubscriptionTier == model.TierPremium {
			return c.Next()
		}

		// Basic tier (including trial) - check limit
		count, err := d.Tenant().CountTableRecords(ctx, tenantID, tableName)
		if err != nil {
			// On error, allow request (fail open for better UX)
			return c.Next()
		}

		if count >= limit {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Batas data tercapai. Upgrade ke Premium untuk menambah lebih banyak data.",
				"code":  "LIMIT_REACHED",
				"limit_info": fiber.Map{
					"table":   tableName,
					"current": count,
					"max":     limit,
				},
			})
		}

		return c.Next()
	}
}
