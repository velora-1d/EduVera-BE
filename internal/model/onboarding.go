package model

import (
	"encoding/json"
	"time"
)

type OnboardingSession struct {
	ID          string          `json:"id" db:"id"`
	TenantID    string          `json:"tenant_id" db:"tenant_id"`
	UserID      string          `json:"user_id" db:"user_id"`
	CurrentStep int             `json:"current_step" db:"current_step"`
	Data        json.RawMessage `json:"data" db:"data"`
	ExpiresAt   *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// OnboardingData stores all data across steps
type OnboardingData struct {
	// Step 1: Account
	AdminName     string `json:"admin_name"`
	AdminEmail    string `json:"admin_email"`
	AdminWhatsApp string `json:"admin_whatsapp"`

	// Step 2: Institution
	InstitutionName string `json:"institution_name"`
	InstitutionType string `json:"institution_type"`
	PlanType        string `json:"plan_type"`
	Address         string `json:"address"`

	// Step 3: Subdomain
	Subdomain string `json:"subdomain"`

	// Step 4: Bank Account
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountHolder string `json:"account_holder"`
}

type OnboardingSessionInput struct {
	TenantID    string          `json:"tenant_id"`
	UserID      string          `json:"user_id"`
	CurrentStep int             `json:"current_step"`
	Data        json.RawMessage `json:"data"`
}

func OnboardingSessionPrepare(input *OnboardingSessionInput) *OnboardingSession {
	expiry := time.Now().Add(24 * time.Hour) // 24 hour expiry

	return &OnboardingSession{
		TenantID:    input.TenantID,
		UserID:      input.UserID,
		CurrentStep: input.CurrentStep,
		Data:        input.Data,
		ExpiresAt:   &expiry,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
