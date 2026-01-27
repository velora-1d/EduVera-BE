package fiber_inbound_adapter

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"eduvera/internal/domain"
	"eduvera/internal/model"
	inbound_port "eduvera/internal/port/inbound"
)

type authAdapter struct {
	domain domain.Domain
}

func NewAuthAdapter(domain domain.Domain) inbound_port.AuthHttpPort {
	return &authAdapter{
		domain: domain,
	}
}

// POST /api/v1/auth/login
func (h *authAdapter) Login(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input model.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	response, err := h.domain.Auth().Login(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":       "success",
		"user":         response.User,
		"access_token": response.AccessToken,
		"expires_at":   response.ExpiresAt,
		"token_type":   "Bearer",
	})
}

// GET /api/v1/auth/me
func (h *authAdapter) Me(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	// Get token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	claims, err := h.domain.Auth().ValidateToken(ctx, tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	user, err := h.domain.Auth().GetCurrentUser(ctx, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// POST /api/v1/auth/refresh - placeholder for token refresh
func (h *authAdapter) Refresh(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Token refresh not yet implemented",
	})
}

// POST /api/v1/auth/logout - placeholder for logout
func (h *authAdapter) Logout(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
