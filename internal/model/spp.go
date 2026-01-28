package model

import (
	"time"
)

type SPPStatus string

const (
	SPPStatusPending SPPStatus = "pending"
	SPPStatusPaid    SPPStatus = "paid"
	SPPStatusOverdue SPPStatus = "overdue"
)

type SPPTransaction struct {
	ID            string    `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	StudentID     string    `json:"student_id,omitempty" db:"student_id"`
	StudentName   string    `json:"student_name" db:"student_name"`
	Amount        int64     `json:"amount" db:"amount"`
	PaymentMethod string    `json:"payment_method,omitempty" db:"payment_method"`
	Status        SPPStatus `json:"status" db:"status"`
	GatewayRef    string    `json:"gateway_ref,omitempty" db:"gateway_ref"`
	Description   string    `json:"description,omitempty" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type SPPStats struct {
	TotalBilled      int64 `json:"total_billed"`
	TotalPaid        int64 `json:"total_paid"`
	TotalPending     int64 `json:"total_pending"`
	TransactionCount int   `json:"transaction_count"`
	PaidCount        int   `json:"paid_count"`
	PendingCount     int   `json:"pending_count"`
}
