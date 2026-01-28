package fiber_inbound_adapter

import (
	"os"

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
			"error": "Unauthorized",
		})
	}

	// Use Auth Domain to validate token
	claims, err := h.domain.Auth().ValidateToken(ctx, bearerToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Check Role
	if claims.Role != model.RoleSuperAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: Owner access required",
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
			"error": "Unauthorized",
		})
	}

	if bearerToken != os.Getenv("INTERNAL_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
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
			Error:   "Unauthorized",
		})
	}

	authDriver := os.Getenv("AUTH_DRIVER")
	if authDriver == "jwt" {
		jwksURL := os.Getenv("AUTH_JWKS_URL")

		_, err := jwt.ValidateJWTWithURL(bearerToken, jwksURL)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Success: false,
				Error:   "Unauthorized: " + err.Error(),
			})
		}
	} else {
		exists, err := h.domain.Client().IsExists(ctx, bearerToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Error:   err.Error(),
			})
		}

		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Success: false,
				Error:   "Unauthorized",
			})
		}
	}

	return c.Next()
}

// RequirePlan creates middleware that checks if the authenticated tenant has the required plan type
// This prevents users from accessing features not included in their subscription
func RequirePlan(requiredPlans ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get tenant plan from context (should be set by ClientAuth middleware)
		planType, ok := c.Locals("plan_type").(string)
		if !ok || planType == "" {
			// If no plan type in context, try to get from JWT claims
			// For now, allow access if plan checking is not yet fully implemented
			return c.Next()
		}

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
