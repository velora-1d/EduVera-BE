package model

import (
	"time"
)

type DisbursementStatus string

const (
	DisbursementStatusPending   DisbursementStatus = "pending"
	DisbursementStatusCompleted DisbursementStatus = "completed"
	DisbursementStatusRejected  DisbursementStatus = "rejected"
)

type Disbursement struct {
	ID            string             `json:"id" db:"id"`
	TenantID      string             `json:"tenant_id" db:"tenant_id"`
	TenantName    string             `json:"tenant_name" db:"tenant_name"` // Joined from tenants table
	Amount        int64              `json:"amount" db:"amount"`
	BankName      string             `json:"bank_name" db:"bank_name"`
	AccountNumber string             `json:"account_number" db:"account_number"`
	AccountHolder string             `json:"account_holder" db:"account_holder"`
	Status        DisbursementStatus `json:"status" db:"status"`
	Notes         string             `json:"notes" db:"notes"`
	AdminNotes    string             `json:"admin_notes" db:"admin_notes"`
	RequestedAt   time.Time          `json:"requested_at" db:"requested_at"`
	ProcessedAt   *time.Time         `json:"processed_at,omitempty" db:"processed_at"`
}
