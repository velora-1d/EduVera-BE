package model

import (
	"time"
)

// Subscription status constants
const (
	SubscriptionStatusActive      = "active"
	SubscriptionStatusGracePeriod = "grace_period"
	SubscriptionStatusSuspended   = "suspended"
	SubscriptionStatusCancelled   = "cancelled"
	SubscriptionStatusTerminated  = "terminated"
)

// Billing cycle constants
const (
	BillingCycleMonthly = "monthly"
	BillingCycleAnnual  = "annual"
)

// Grace period duration in days
const (
	GracePeriodDays = 7
	SuspensionDays  = 90
	TerminationDays = 90
)

// Subscription represents tenant subscription
type Subscription struct {
	ID                 string     `json:"id" db:"id"`
	TenantID           string     `json:"tenant_id" db:"tenant_id"`
	PlanType           string     `json:"plan_type" db:"plan_type"`
	SubscriptionTier   string     `json:"subscription_tier" db:"subscription_tier"` // basic or premium
	BillingCycle       string     `json:"billing_cycle" db:"billing_cycle"`
	Status             string     `json:"status" db:"status"`
	CurrentPeriodStart time.Time  `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end" db:"current_period_end"`
	GracePeriodEnd     time.Time  `json:"grace_period_end" db:"grace_period_end"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	ScheduledPlanType  *string    `json:"scheduled_plan_type,omitempty" db:"scheduled_plan_type"` // For downgrade
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// PricingPlan represents pricing configuration from database
type PricingPlan struct {
	ID           string    `json:"id" db:"id"`
	PlanType     string    `json:"plan_type" db:"plan_type"`
	BillingCycle string    `json:"billing_cycle" db:"billing_cycle"`
	Price        int64     `json:"price" db:"price"`
	Currency     string    `json:"currency" db:"currency"`
	Description  string    `json:"description" db:"description"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSubscriptionInput for creating new subscription
type CreateSubscriptionInput struct {
	TenantID     string `json:"tenant_id" validate:"required"`
	PlanType     string `json:"plan_type" validate:"required,oneof=sekolah pesantren hybrid"`
	BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly annual"`
}

// UpgradeInput for upgrading plan
type UpgradeInput struct {
	TenantID    string `json:"tenant_id" validate:"required"`
	NewPlanType string `json:"new_plan_type" validate:"required,oneof=sekolah pesantren hybrid"`
}

// DowngradeInput for downgrading plan
type DowngradeInput struct {
	TenantID    string `json:"tenant_id" validate:"required"`
	NewPlanType string `json:"new_plan_type" validate:"required,oneof=sekolah pesantren hybrid"`
	Confirmed   bool   `json:"confirmed"` // User confirmed data will be hidden
}

// SubscriptionFilter for querying subscriptions
type SubscriptionFilter struct {
	IDs       []string `json:"ids"`
	TenantIDs []string `json:"tenant_ids"`
	Statuses  []string `json:"statuses"`
	PlanTypes []string `json:"plan_types"`
}

// PricingFilter for querying pricing plans
type PricingFilter struct {
	PlanTypes     []string `json:"plan_types"`
	BillingCycles []string `json:"billing_cycles"`
	IsActive      *bool    `json:"is_active"`
}

// UpgradeCalculation represents prorata calculation result
type UpgradeCalculation struct {
	OldPlanType   string `json:"old_plan_type"`
	NewPlanType   string `json:"new_plan_type"`
	RemainingDays int    `json:"remaining_days"`
	OldPrice      int64  `json:"old_price"`
	NewPrice      int64  `json:"new_price"`
	ProrataCredit int64  `json:"prorata_credit"`
	AmountDue     int64  `json:"amount_due"`
}

// DowngradeValidation result
type DowngradeValidation struct {
	CanDowngrade    bool           `json:"can_downgrade"`
	AffectedModules []string       `json:"affected_modules"`
	DataCount       map[string]int `json:"data_count"`
	EffectiveDate   time.Time      `json:"effective_date"`
	Warning         string         `json:"warning,omitempty"`
}

// IsExpired checks if subscription is past current period
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.CurrentPeriodEnd)
}

// IsInGracePeriod checks if in grace period
func (s *Subscription) IsInGracePeriod() bool {
	now := time.Now()
	return now.After(s.CurrentPeriodEnd) && now.Before(s.GracePeriodEnd)
}

// DaysRemaining returns remaining days in current period
func (s *Subscription) DaysRemaining() int {
	remaining := time.Until(s.CurrentPeriodEnd)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Hours() / 24)
}

// TotalDays returns total days in billing cycle
func (s *Subscription) TotalDays() int {
	if s.BillingCycle == BillingCycleAnnual {
		return 365
	}
	return 30
}

func (f SubscriptionFilter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.TenantIDs) == 0 && len(f.Statuses) == 0 && len(f.PlanTypes) == 0
}

func (f PricingFilter) IsEmpty() bool {
	return len(f.PlanTypes) == 0 && len(f.BillingCycles) == 0 && f.IsActive == nil
}
