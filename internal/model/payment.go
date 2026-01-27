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
var PlanPricing = map[string]map[bool]int64{
	"sekolah": {
		true:  4990000, // Annual
		false: 499000,  // Monthly
	},
	"pesantren": {
		true:  4990000, // Annual
		false: 499000,  // Monthly
	},
	"hybrid": {
		true:  7990000, // Annual
		false: 799000,  // Monthly
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
	PlanType string `json:"plan_type" validate:"required,oneof=sekolah pesantren hybrid"`
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

// GetPlanPrice returns price based on plan and billing cycle
func GetPlanPrice(planType string, isAnnual bool) int64 {
	if prices, ok := PlanPricing[planType]; ok {
		return prices[isAnnual]
	}
	return 0
}

func (f PaymentFilter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.TenantIDs) == 0 && len(f.OrderIDs) == 0 && len(f.Statuses) == 0
}
