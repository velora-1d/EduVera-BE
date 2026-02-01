package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// ResetToken represents a password reset token
type ResetToken struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// ResetTokenExpiry is the duration a reset token is valid
const ResetTokenExpiry = 24 * time.Hour

// IsExpired checks if the token has expired
func (t *ResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *ResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (t *ResetToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// GenerateResetToken generates a secure random token
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ForgotPasswordInput for requesting password reset
type ForgotPasswordInput struct {
	Email    string `json:"email" validate:"required,email"`
	TenantID string `json:"tenant_id" validate:"required"`
}

// ResetPasswordInput for resetting password with token
type ResetPasswordInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
