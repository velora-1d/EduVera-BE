package rabbitmq_outbound_adapter

import (
	"context"
	"prabogo/utils/rabbitmq"
)

const (
	NotificationTypeSystem = "system"
	WhatsAppExchange       = "whatsapp_notifications"
)

type WhatsAppNotification struct {
	Type      string `json:"type"`
	TenantID  string `json:"tenant_id,omitempty"`
	Phone     string `json:"phone"`
	Message   string `json:"message"`
	Priority  int    `json:"priority"`
	RetryLeft int    `json:"retry_left"`
}

type whatsappAdapter struct{}

func NewWhatsAppAdapter() *whatsappAdapter {
	return &whatsappAdapter{}
}

func (w *whatsappAdapter) Send(target, message string) error {
	notif := WhatsAppNotification{
		Type:      NotificationTypeSystem,
		Phone:     target,
		Message:   message,
		RetryLeft: 3,
	}

	return rabbitmq.Publish(context.Background(), WhatsAppExchange, rabbitmq.KindFanOut, "", notif)
}
