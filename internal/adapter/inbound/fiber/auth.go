package fiber_inbound_adapter

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
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
			"error": "Data tidak valid. Silakan coba lagi.",
		})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email dan password wajib diisi.",
		})
	}

	response, err := h.domain.Auth().Login(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email atau password salah. Silakan coba lagi.",
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
			"error": "Sesi Anda telah berakhir. Silakan login kembali.",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Format authorization tidak valid.",
		})
	}

	claims, err := h.domain.Auth().ValidateToken(ctx, tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Sesi Anda telah berakhir. Silakan login kembali.",
		})
	}

	user, err := h.domain.Auth().GetCurrentUser(ctx, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data pengguna tidak ditemukan.",
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
