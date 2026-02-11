package fiber_inbound_adapter

import (
	"context"
	"strings"

	"os"
	"time"

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
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	if input.Email == "" || input.Password == "" {
		return SendError(c, fiber.StatusBadRequest, "Email dan password wajib diisi.", nil)
	}

	response, err := h.domain.Auth().Login(ctx, &input, c.IP())
	if err != nil {
		// Log actual error but show generic message
		return SendError(c, fiber.StatusUnauthorized, "Email atau password salah. Silakan coba lagi.", err)
	}

	// Set HttpOnly Cookie
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ".ve-lora.my.id" // Fallback
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    response.AccessToken,
		Expires:  time.Unix(response.ExpiresAt, 0), // response.ExpiresAt is likely int64 timestamp
		HTTPOnly: true,
		Secure:   true,         // Always true for "None" SameSite
		SameSite: "None",       // Required for cross-site (api -> app)
		Domain:   cookieDomain, // Cross-subdomain
		Path:     "/",
	})

	return SendSuccess(c, "Login berhasil", fiber.Map{
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
		return SendError(c, fiber.StatusUnauthorized, "Sesi Anda telah berakhir. Silakan login kembali.", nil)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return SendError(c, fiber.StatusUnauthorized, "Format authorization tidak valid.", nil)
	}

	claims, err := h.domain.Auth().ValidateToken(ctx, tokenString)
	if err != nil {
		return SendError(c, fiber.StatusUnauthorized, "Sesi Anda telah berakhir. Silakan login kembali.", err)
	}

	user, err := h.domain.Auth().GetCurrentUser(ctx, claims.UserID)
	if err != nil {
		return SendError(c, fiber.StatusNotFound, "Data pengguna tidak ditemukan.", err)
	}

	return SendSuccess(c, "Data user berhasil diambil", fiber.Map{
		"user": user,
	})
}

// POST /api/v1/auth/refresh - placeholder for token refresh
func (h *authAdapter) Refresh(a any) error {
	c := a.(*fiber.Ctx)
	return SendSuccess(c, "Token refresh not yet implemented", nil)
}

// POST /api/v1/auth/logout - blacklist token to invalidate session
func (h *authAdapter) Logout(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	// Get token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Berhasil logout",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader || tokenString == "" {
		return SendSuccess(c, "Berhasil logout", nil)
	}

	// Validate token to get expiry time
	claims, err := h.domain.Auth().ValidateToken(ctx, tokenString)
	if err == nil && claims != nil {
		// Blacklist the token until its expiry
		expiresAt := claims.ExpiresAt.Time
		_ = h.domain.Auth().BlacklistToken(ctx, tokenString, expiresAt)
	}

	// Clear Cookie
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ".ve-lora.my.id"
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Domain:   cookieDomain,
		Path:     "/",
	})

	return SendSuccess(c, "Berhasil logout", nil)
}

// POST /api/v1/auth/forgot-password
func (h *authAdapter) ForgotPassword(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input model.ForgotPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	if input.Email == "" {
		return SendError(c, fiber.StatusBadRequest, "Email wajib diisi.", nil)
	}

	// Process forgot password (always return success to not reveal email existence)
	_ = h.domain.Auth().ForgotPassword(ctx, &input)

	return SendSuccess(c, "Jika email terdaftar, link reset password akan dikirim ke WhatsApp Anda.", nil)
}

// POST /api/v1/auth/reset-password
func (h *authAdapter) ResetPassword(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	var input model.ResetPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Data tidak valid. Silakan coba lagi.", err)
	}

	if input.Token == "" || input.NewPassword == "" {
		return SendError(c, fiber.StatusBadRequest, "Token dan password baru wajib diisi.", nil)
	}

	if len(input.NewPassword) < 8 {
		return SendError(c, fiber.StatusBadRequest, "Password minimal 8 karakter.", nil)
	}

	err := h.domain.Auth().ResetPassword(ctx, &input)
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, "Token tidak valid atau sudah kadaluarsa. Silakan request reset password baru.", err)
	}

	return SendSuccess(c, "Password berhasil diubah. Silakan login dengan password baru Anda.", nil)
}
