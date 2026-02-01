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

	// Step 2: Institution & Plan
	InstitutionName  string   `json:"institution_name"`          // Primary name / Yayasan name
	SchoolName       string   `json:"school_name,omitempty"`     // For hybrid: specific school name
	PesantrenName    string   `json:"pesantren_name,omitempty"`  // For hybrid: specific pesantren name
	SchoolJenjangs   []string `json:"school_jenjangs,omitempty"` // Multi-select: TK, SD, MI, SMP, MTs, SMA, MA, SMK
	InstitutionType  string   `json:"institution_type"`          // sekolah, pesantren, hybrid
	PlanType         string   `json:"plan_type"`                 // sekolah, pesantren, hybrid
	SubscriptionTier string   `json:"subscription_tier"`         // basic, premium
	BillingCycle     string   `json:"billing_cycle"`             // monthly, annual
	Address          string   `json:"address"`

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
