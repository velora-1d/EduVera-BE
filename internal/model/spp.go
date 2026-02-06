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

// PaymentType differentiates between school (SPP) and pesantren (Syahriah) payments
type PaymentType string

const (
	PaymentTypeSPP      PaymentType = "spp"      // Sekolah - school tuition
	PaymentTypeSyahriah PaymentType = "syahriah" // Pesantren - boarding fee
	PaymentTypeOther    PaymentType = "other"    // Other fees (registration, books, etc)
)

type SPPTransaction struct {
	ID            string      `json:"id" db:"id"`
	TenantID      string      `json:"tenant_id" db:"tenant_id"`
	StudentID     string      `json:"student_id,omitempty" db:"student_id"`
	StudentName   string      `json:"student_name" db:"student_name"`
	Amount        int64       `json:"amount" db:"amount"`
	PaymentMethod string      `json:"payment_method,omitempty" db:"payment_method"` // cash, transfer, qris
	PaymentType   PaymentType `json:"payment_type" db:"payment_type"`               // spp, syahriah, other
	Status        SPPStatus   `json:"status" db:"status"`
	GatewayRef    string      `json:"gateway_ref,omitempty" db:"gateway_ref"`
	Description   string      `json:"description,omitempty" db:"description"`
	PaymentProof  string      `json:"payment_proof,omitempty" db:"payment_proof"` // URL bukti transfer
	ConfirmedBy   string      `json:"confirmed_by,omitempty" db:"confirmed_by"`   // User ID yang konfirmasi
	PaidAt        *time.Time  `json:"paid_at,omitempty" db:"paid_at"`             // Waktu pembayaran dikonfirmasi
	DueDate       *time.Time  `json:"due_date,omitempty" db:"due_date"`           // Tanggal jatuh tempo
	Period        string      `json:"period,omitempty" db:"period"`               // Periode: "2024-01" for Jan 2024
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

type SPPStats struct {
	TotalBilled      int64 `json:"total_billed"`
	TotalPaid        int64 `json:"total_paid"`
	TotalPending     int64 `json:"total_pending"`
	TransactionCount int   `json:"transaction_count"`
	PaidCount        int   `json:"paid_count"`
	PendingCount     int   `json:"pending_count"`
}
