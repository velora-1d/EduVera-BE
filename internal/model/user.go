package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User roles
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleTeacher    = "teacher"
	RoleStudent    = "student"
	RoleParent     = "parent"
)

type User struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	Name            string     `json:"name" db:"name"`
	Email           string     `json:"email" db:"email"`
	WhatsApp        string     `json:"whatsapp,omitempty" db:"whatsapp"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	Role            string     `json:"role" db:"role"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type UserInput struct {
	TenantID string `json:"tenant_id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	WhatsApp string `json:"whatsapp,omitempty"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role,omitempty"`
}

type UserFilter struct {
	IDs       []string `json:"ids"`
	TenantIDs []string `json:"tenant_ids"`
	Emails    []string `json:"emails"`
	Roles     []string `json:"roles"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func UserPrepare(input *UserInput) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := input.Role
	if role == "" {
		role = RoleAdmin
	}

	return &User{
		TenantID:     input.TenantID,
		Name:         input.Name,
		Email:        input.Email,
		WhatsApp:     input.WhatsApp,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (f UserFilter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.TenantIDs) == 0 && len(f.Emails) == 0 && len(f.Roles) == 0
}
