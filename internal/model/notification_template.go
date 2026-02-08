package model

import "time"

// NotificationTemplate represents a customizable notification template
type NotificationTemplate struct {
	ID              string    `json:"id" db:"id"`
	EventType       string    `json:"event_type" db:"event_type"`
	Channel         string    `json:"channel" db:"channel"`
	TemplateName    string    `json:"template_name" db:"template_name"`
	TemplateContent string    `json:"template_content" db:"template_content"`
	Variables       string    `json:"variables" db:"variables"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Event types
const (
	NotificationEventRegistration = "registration"
	NotificationEventUpgrade      = "upgrade"
	NotificationEventPayment      = "payment"
)

// Channel types
const (
	NotificationChannelTelegram       = "telegram"
	NotificationChannelWhatsAppOwner  = "whatsapp_owner"
	NotificationChannelWhatsAppTenant = "whatsapp_tenant"
)
