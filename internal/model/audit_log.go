package model

import "time"

// AuditAction types
const (
	// Existing actions
	AuditActionTenantStatusUpdate  = "tenant_status_update"
	AuditActionDisbursementApprove = "disbursement_approve"
	AuditActionDisbursementReject  = "disbursement_reject"
	AuditActionContentUpdate       = "content_update"
	AuditActionUserBan             = "user_ban"

	// Subscription actions
	AuditActionSubscriptionCreated   = "subscription_created"
	AuditActionSubscriptionUpgraded  = "subscription_upgraded"
	AuditActionSubscriptionRenewed   = "subscription_renewed"
	AuditActionSubscriptionExpired   = "subscription_expired"
	AuditActionSubscriptionSuspended = "subscription_suspended"

	// Payment actions
	AuditActionPaymentSuccess = "payment_success"
	AuditActionPaymentFailed  = "payment_failed"
	AuditActionPaymentRefund  = "payment_refund"

	// Auth actions
	AuditActionLoginSuccess = "login_success"
	AuditActionLoginFailed  = "login_failed"
	AuditActionLogout       = "logout"

	// WhatsApp actions
	AuditActionWhatsAppConnected    = "whatsapp_connected"
	AuditActionWhatsAppDisconnected = "whatsapp_disconnected"
)

// AuditLog represents an admin action log entry
type AuditLog struct {
	ID          string    `json:"id" db:"id"`
	AdminID     string    `json:"admin_id" db:"admin_id"`
	AdminEmail  string    `json:"admin_email" db:"admin_email"`
	Action      string    `json:"action" db:"action"`
	TargetType  string    `json:"target_type" db:"target_type"` // tenant, disbursement, user, content
	TargetID    string    `json:"target_id" db:"target_id"`
	OldValue    string    `json:"old_value,omitempty" db:"old_value"`
	NewValue    string    `json:"new_value,omitempty" db:"new_value"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AuditLogInput for creating new audit log
type AuditLogInput struct {
	AdminID     string
	AdminEmail  string
	Action      string
	TargetType  string
	TargetID    string
	OldValue    string
	NewValue    string
	IPAddress   string
	UserAgent   string
	Description string
}

// AuditLogFilter for querying logs
type AuditLogFilter struct {
	AdminID    string
	Action     string
	TargetType string
	TargetID   string
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}
