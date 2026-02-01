package model

import (
	"fmt"
	"time"
)

// Payment status constants
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusFailed   = "failed"
	PaymentStatusExpired  = "expired"
	PaymentStatusRefunded = "refunded"
)

// Plan pricing in Rupiah
// Structure: plan_tier -> isAnnual -> price
// Plans: sekolah_basic, pesantren_basic, hybrid_basic, sekolah_premium, pesantren_premium, hybrid_premium
var PlanPricing = map[string]map[bool]int64{
	// BASIC PLANS
	"sekolah_basic": {
		true:  2999000, // Annual - Rp 2.999.000
		false: 299000,  // Monthly - Rp 299.000
	},
	"pesantren_basic": {
		true:  2999000, // Annual - Rp 2.999.000
		false: 299000,  // Monthly - Rp 299.000
	},
	"hybrid_basic": {
		true:  4499000, // Annual - Rp 4.499.000
		false: 449000,  // Monthly - Rp 449.000
	},
	// PREMIUM PLANS
	"sekolah_premium": {
		true:  4999000, // Annual - Rp 4.999.000
		false: 499000,  // Monthly - Rp 499.000
	},
	"pesantren_premium": {
		true:  4999000, // Annual - Rp 4.999.000
		false: 499000,  // Monthly - Rp 499.000
	},
	"hybrid_premium": {
		true:  6999999, // Annual - Rp 6.999.999 (FLAGSHIP)
		false: 699000,  // Monthly - Rp 699.000
	},
	// Legacy support - map old plan names to premium
	"sekolah": {
		true:  4999000,
		false: 499000,
	},
	"pesantren": {
		true:  4999000,
		false: 499000,
	},
	"hybrid": {
		true:  6999999,
		false: 699000,
	},
}

type Payment struct {
	ID          string     `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	OrderID     string     `json:"order_id" db:"order_id"`
	Amount      int64      `json:"amount" db:"amount"`
	Status      string     `json:"status" db:"status"`
	PaymentType string     `json:"payment_type,omitempty" db:"payment_type"`
	SnapToken   string     `json:"snap_token,omitempty" db:"snap_token"`
	MidtransID  string     `json:"midtrans_id,omitempty" db:"midtrans_id"`
	PaidAt      *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreatePaymentInput struct {
	TenantID string `json:"tenant_id" validate:"required"`
	PlanType string `json:"plan_type" validate:"required"`
	PlanTier string `json:"plan_tier"` // basic or premium
	IsAnnual bool   `json:"is_annual"`
}

type PaymentFilter struct {
	IDs       []string `json:"ids"`
	TenantIDs []string `json:"tenant_ids"`
	OrderIDs  []string `json:"order_ids"`
	Statuses  []string `json:"statuses"`
}

// SnapTransactionResponse from Midtrans
type SnapTransactionResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// MidtransNotification webhook payload
type MidtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

// GenerateOrderID creates unique order ID
func GenerateOrderID(tenantID string) string {
	timestamp := time.Now().UnixMilli()
	return fmt.Sprintf("EDV-%d-%s", timestamp, tenantID[:8])
}

// GenerateTimestamp returns current timestamp as string for order IDs
func GenerateTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

// GetPlanPrice returns price based on plan type, tier and billing cycle
func GetPlanPrice(planType string, isAnnual bool) int64 {
	if prices, ok := PlanPricing[planType]; ok {
		return prices[isAnnual]
	}
	return 0
}

// GetPlanPriceWithTier returns price with explicit tier
func GetPlanPriceWithTier(planType, tier string, isAnnual bool) int64 {
	key := planType + "_" + tier
	if prices, ok := PlanPricing[key]; ok {
		return prices[isAnnual]
	}
	// Fallback to legacy
	return GetPlanPrice(planType, isAnnual)
}

func (f PaymentFilter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.TenantIDs) == 0 && len(f.OrderIDs) == 0 && len(f.Statuses) == 0
}
