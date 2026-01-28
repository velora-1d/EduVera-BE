package model

import (
	"time"
)

// Notification types
const (
	NotificationTypeTelegram = "telegram"
	NotificationTypeEmail    = "email"
	NotificationTypeWhatsapp = "whatsapp"
)

// Notification status
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

type Notification struct {
	ID        string    `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`           // telegram, email, whatsapp
	Recipient string    `json:"recipient" db:"recipient"` // channel/email/phone
	Subject   string    `json:"subject,omitempty" db:"subject"`
	Message   string    `json:"message" db:"message"`
	Status    string    `json:"status" db:"status"`
	TenantID  string    `json:"tenant_id,omitempty" db:"tenant_id"`
	ErrorMsg  string    `json:"error_msg,omitempty" db:"error_msg"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type NotificationInput struct {
	Type      string `json:"type" validate:"required,oneof=telegram email whatsapp"`
	Recipient string `json:"recipient" validate:"required"`
	Message   string `json:"message" validate:"required"`
	TenantID  string `json:"tenant_id,omitempty"`
}

type NotificationFilter struct {
	Types    []string `json:"types"`
	Statuses []string `json:"statuses"`
	TenantID string   `json:"tenant_id"`
}

type NotificationStats struct {
	TotalSent   int `json:"total_sent"`
	TotalFailed int `json:"total_failed"`
}

func NotificationPrepare(input *NotificationInput) *Notification {
	return &Notification{
		Type:      input.Type,
		Recipient: input.Recipient,
		Message:   input.Message,
		TenantID:  input.TenantID,
		Status:    NotificationStatusPending,
		CreatedAt: time.Now(),
	}
}

func (f NotificationFilter) IsEmpty() bool {
	return len(f.Types) == 0 && len(f.Statuses) == 0 && f.TenantID == ""
}
