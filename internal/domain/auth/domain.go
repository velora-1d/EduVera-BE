package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/palantir/stacktrace"
	"golang.org/x/crypto/bcrypt"

	"prabogo/internal/domain/audit_log"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
	"prabogo/utils/redis"
)

type AuthDomain interface {
	Register(ctx context.Context, input *model.UserInput) (*model.User, error)
	Login(ctx context.Context, input *model.LoginInput, ipAddress string) (*model.LoginResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
	GetCurrentUser(ctx context.Context, userID string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GenerateToken(user *model.User) (string, int64, error)
	LinkUserToTenant(ctx context.Context, userID string, tenantID string) error
	ForgotPassword(ctx context.Context, input *model.ForgotPasswordInput) error
	ResetPassword(ctx context.Context, input *model.ResetPasswordInput) error
	// Token blacklist for secure logout
	BlacklistToken(ctx context.Context, token string, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type authDomain struct {
	databasePort outbound_port.DatabasePort
	messagePort  outbound_port.MessagePort
}

func NewAuthDomain(databasePort outbound_port.DatabasePort, messagePort outbound_port.MessagePort) AuthDomain {
	return &authDomain{
		databasePort: databasePort,
		messagePort:  messagePort,
	}
}

func (d *authDomain) Register(ctx context.Context, input *model.UserInput) (*model.User, error) {
	// Check if email already exists
	exists, err := d.databasePort.User().EmailExists(input.Email)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to check email")
	}
	if exists {
		return nil, stacktrace.NewError("registration failed, please try again")
	}

	user, err := model.UserPrepare(input)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to prepare user")
	}

	err = d.databasePort.User().Create(user)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to create user")
	}

	// Send Welcome WhatsApp Notification
	if d.messagePort != nil && user.WhatsApp != "" {
		message := "Halo " + user.Name + "!\n\nSelamat datang di EduVera. Akun admin Anda telah berhasil dibuat.\n\n" +
			"Email: " + user.Email + "\n" +
			"Silakan lanjutkan proses pendaftaran dengan melengkapi data institusi Anda.\n\n" +
			"Terima kasih!"
		_ = d.messagePort.WhatsApp().Send(user.WhatsApp, message)
	}

	return user, nil
}

func (d *authDomain) Login(ctx context.Context, input *model.LoginInput, ipAddress string) (*model.LoginResponse, error) {
	audit := audit_log.NewAuditHelper(d.databasePort.AuditLog())

	user, err := d.databasePort.User().FindByEmail(input.Email)
	if err != nil {
		_ = audit.LogLoginEvent(ctx, model.AuditActionLoginFailed, "", input.Email, ipAddress)
		return nil, stacktrace.NewError("invalid email or password")
	}

	if !user.CheckPassword(input.Password) {
		_ = audit.LogLoginEvent(ctx, model.AuditActionLoginFailed, user.ID, user.Email, ipAddress)
		return nil, stacktrace.NewError("invalid email or password")
	}

	if !user.IsActive {
		_ = audit.LogLoginEvent(ctx, model.AuditActionLoginFailed, user.ID, user.Email, ipAddress)
		return nil, stacktrace.NewError("account not activated")
	}

	// Update last login
	_ = d.databasePort.User().UpdateLastLogin(user.ID)

	// Log Success
	_ = audit.LogLoginEvent(ctx, model.AuditActionLoginSuccess, user.ID, user.Email, ipAddress)

	// Generate JWT
	signedToken, expiresAt, err := d.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		User:        *user,
		AccessToken: signedToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (d *authDomain) GenerateToken(user *model.User) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eduvera",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return "", 0, stacktrace.Propagate(err, "failed to sign token")
	}

	return signedToken, expiresAt.Unix(), nil
}

func (d *authDomain) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, stacktrace.NewError("unexpected signing method")
		}
		return []byte(getJWTSecret()), nil
	})

	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to parse token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, stacktrace.NewError("invalid token")
	}

	return claims, nil
}

func (d *authDomain) GetCurrentUser(ctx context.Context, userID string) (*model.User, error) {
	user, err := d.databasePort.User().FindByID(userID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find user")
	}
	return user, nil
}

func (d *authDomain) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := d.databasePort.User().FindByEmail(email)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find user by email")
	}
	return user, nil
}

func (d *authDomain) LinkUserToTenant(ctx context.Context, userID string, tenantID string) error {
	// Optional: Check if user exists
	// Optional: Check if tenant exists (usually ensured by caller)

	err := d.databasePort.User().LinkToTenant(userID, tenantID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to link user to tenant")
	}

	return nil
}

func (d *authDomain) ForgotPassword(ctx context.Context, input *model.ForgotPasswordInput) error {
	// Find user by email
	user, err := d.databasePort.User().FindByEmail(input.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate reset token
	tokenStr, err := model.GenerateResetToken()
	if err != nil {
		return stacktrace.Propagate(err, "failed to generate reset token")
	}

	// Create reset token record
	resetToken := &model.ResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(model.ResetTokenExpiry),
		CreatedAt: time.Now(),
	}

	err = d.databasePort.User().CreateResetToken(resetToken)
	if err != nil {
		return stacktrace.Propagate(err, "failed to create reset token")
	}

	// Build reset link
	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "https://eduvera.ve-lora.my.id"
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, tokenStr)

	// Send WhatsApp notification
	if d.messagePort != nil && user.WhatsApp != "" {
		message := fmt.Sprintf(
			"🔐 *Reset Password EduVera*\n\n"+
				"Halo %s,\n\n"+
				"Kami menerima permintaan untuk reset password akun Anda.\n\n"+
				"Klik link berikut untuk membuat password baru:\n%s\n\n"+
				"Link ini berlaku selama 24 jam.\n\n"+
				"Jika Anda tidak meminta reset password, abaikan pesan ini.\n\n"+
				"Terima kasih,\nTim EduVera",
			user.Name, resetLink,
		)
		_ = d.messagePort.WhatsApp().Send(user.WhatsApp, message)
	}

	return nil
}

func (d *authDomain) ResetPassword(ctx context.Context, input *model.ResetPasswordInput) error {
	// Get reset token
	resetToken, err := d.databasePort.User().GetResetToken(input.Token)
	if err != nil {
		return stacktrace.NewError("token tidak valid atau sudah kadaluarsa")
	}

	// Validate token
	if !resetToken.IsValid() {
		return stacktrace.NewError("token tidak valid atau sudah digunakan")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return stacktrace.Propagate(err, "failed to hash password")
	}

	// Update password
	err = d.databasePort.User().UpdatePassword(resetToken.UserID, string(hashedPassword))
	if err != nil {
		return stacktrace.Propagate(err, "failed to update password")
	}

	// Mark token as used
	err = d.databasePort.User().MarkResetTokenUsed(resetToken.ID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to mark token as used")
	}

	return nil
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("CRITICAL: JWT_SECRET environment variable is required but not set")
	}
	if len(secret) < 32 {
		panic("CRITICAL: JWT_SECRET must be at least 32 characters long")
	}
	return secret
}

// Token blacklist prefix for Redis keys
const tokenBlacklistPrefix = "bl:"

// hashToken creates a short hash of the token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:8]) // Use first 8 bytes = 16 hex chars
}

// BlacklistToken adds a token to the blacklist until its expiry
func (d *authDomain) BlacklistToken(ctx context.Context, token string, expiresAt time.Time) error {
	key := tokenBlacklistPrefix + hashToken(token)
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}
	return redis.SetWithTTL(ctx, key, "1", ttl)
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (d *authDomain) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := tokenBlacklistPrefix + hashToken(token)
	return redis.Exists(ctx, key)
}
