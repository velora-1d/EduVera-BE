package notification

import (
	"context"
	"errors"
	"log"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

// NotificationType determines which provider to use
type NotificationType string

const (
	NotificationTypeSystem NotificationType = "system" // Uses Fonnte (platform-level)
	NotificationTypeTenant NotificationType = "tenant" // Uses Evolution API (tenant-level)
)

// WhatsAppNotification represents a notification to be sent
type WhatsAppNotification struct {
	Type      NotificationType `json:"type"`
	TenantID  string           `json:"tenant_id,omitempty"` // Required for tenant type
	Phone     string           `json:"phone"`
	Message   string           `json:"message"`
	Priority  int              `json:"priority"` // 1=high, 2=normal, 3=low
	RetryLeft int              `json:"retry_left"`
}

// NotificationService handles sending WhatsApp notifications
type NotificationService struct {
	fonntePort    outbound_port.WhatsAppMessagePort
	evolutionPort outbound_port.EvolutionApiPort
	dbPort        outbound_port.DatabasePort
}

func NewNotificationService(
	fonntePort outbound_port.WhatsAppMessagePort,
	evolutionPort outbound_port.EvolutionApiPort,
	dbPort outbound_port.DatabasePort,
) *NotificationService {
	return &NotificationService{
		fonntePort:    fonntePort,
		evolutionPort: evolutionPort,
		dbPort:        dbPort,
	}
}

// Send routes the notification to appropriate provider
func (s *NotificationService) Send(ctx context.Context, notif WhatsAppNotification) error {
	switch notif.Type {
	case NotificationTypeSystem:
		return s.sendViaFonnte(notif)
	case NotificationTypeTenant:
		return s.sendViaTenant(ctx, notif)
	default:
		return errors.New("unknown notification type")
	}
}

// sendViaFonnte sends using platform Fonnte account
func (s *NotificationService) sendViaFonnte(notif WhatsAppNotification) error {
	if s.fonntePort == nil {
		return errors.New("fonnte adapter not configured")
	}
	return s.fonntePort.Send(notif.Phone, notif.Message)
}

// sendViaTenant sends using tenant's connected WhatsApp via Evolution API
func (s *NotificationService) sendViaTenant(ctx context.Context, notif WhatsAppNotification) error {
	if notif.TenantID == "" {
		return errors.New("tenant_id required for tenant notifications")
	}

	// Get tenant's WhatsApp session
	session, err := s.dbPort.WhatsApp().GetByTenantID(ctx, notif.TenantID)
	if err != nil || session == nil {
		log.Printf("[WARN] Tenant %s has no WhatsApp session, falling back to Fonnte", notif.TenantID)
		return s.sendViaFonnte(notif) // Fallback to Fonnte if tenant not connected
	}

	if session.Status != model.WhatsAppStatusConnected {
		log.Printf("[WARN] Tenant %s WhatsApp disconnected, falling back to Fonnte", notif.TenantID)
		return s.sendViaFonnte(notif) // Fallback
	}

	// Send via Evolution API
	return s.evolutionPort.SendMessage(ctx, session.InstanceName, session.APIKey, notif.Phone, notif.Message)
}

// SendBatch sends multiple notifications with rate limiting
func (s *NotificationService) SendBatch(ctx context.Context, notifications []WhatsAppNotification) []error {
	errors := make([]error, len(notifications))

	for i, notif := range notifications {
		err := s.Send(ctx, notif)
		errors[i] = err

		if err != nil {
			log.Printf("[ERROR] Failed to send notification to %s: %v", notif.Phone, err)
		}
	}

	return errors
}
