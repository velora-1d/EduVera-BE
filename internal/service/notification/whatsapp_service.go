package notification

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

// NotificationType determines which provider to use
type NotificationType string

const (
	NotificationTypeSystem NotificationType = "system" // Uses Fonnte (platform-level fallback)
	NotificationTypeTenant NotificationType = "tenant" // Uses Evolution API (tenant-level)
	NotificationTypeOwner  NotificationType = "owner"  // Uses Evolution API (owner's connected WA)
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
	case NotificationTypeOwner:
		return s.sendViaOwner(ctx, notif)
	default:
		return errors.New("unknown notification type")
	}
}

// sendViaFonnte sends using platform Fonnte account (fallback)
func (s *NotificationService) sendViaFonnte(notif WhatsAppNotification) error {
	if s.fonntePort == nil {
		return errors.New("fonnte adapter not configured")
	}
	return s.fonntePort.Send(notif.Phone, notif.Message)
}

// sendViaOwner sends using owner's connected WhatsApp via Evolution API
func (s *NotificationService) sendViaOwner(ctx context.Context, notif WhatsAppNotification) error {
	// Get owner's WhatsApp session (empty tenantID = owner)
	session, err := s.dbPort.WhatsApp().GetByInstanceName(ctx, "eduvera_owner")
	if err != nil || session == nil {
		log.Printf("[WARN] Owner WhatsApp not connected, falling back to Fonnte")
		return s.sendViaFonnte(notif)
	}

	if session.Status != model.WhatsAppStatusConnected {
		log.Printf("[WARN] Owner WhatsApp disconnected, falling back to Fonnte")
		return s.sendViaFonnte(notif)
	}

	// Send via Evolution API
	return s.evolutionPort.SendMessage(ctx, session.InstanceName, session.APIKey, notif.Phone, notif.Message)
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
	errs := make([]error, len(notifications))

	for i, notif := range notifications {
		err := s.Send(ctx, notif)
		errs[i] = err

		if err != nil {
			log.Printf("[ERROR] Failed to send notification to %s: %v", notif.Phone, err)
		}
	}

	return errs
}

// SendMultiChannel sends notification to all relevant channels based on event type
func (s *NotificationService) SendMultiChannel(ctx context.Context, eventType string, variables map[string]string, recipientPhone string, tenantID string) error {
	// Get templates for this event
	templates, err := s.dbPort.NotificationTemplate().GetActiveByEvent(ctx, eventType)
	if err != nil {
		log.Printf("[ERROR] Failed to get templates for event %s: %v", eventType, err)
		return err
	}

	ownerPhone := os.Getenv("OWNER_PHONE")
	if ownerPhone == "" {
		ownerPhone = "6285117776596" // Default owner phone
	}

	for _, template := range templates {
		message := s.renderTemplate(template.TemplateContent, variables)

		switch template.Channel {
		case model.NotificationChannelWhatsAppOwner:
			// Send to owner via owner's connected WA
			err := s.Send(ctx, WhatsAppNotification{
				Type:    NotificationTypeOwner,
				Phone:   ownerPhone,
				Message: message,
			})
			if err != nil {
				log.Printf("[ERROR] Failed to send owner WA notification: %v", err)
			}

		case model.NotificationChannelWhatsAppTenant:
			// Send to recipient via tenant's WA (or fallback)
			if recipientPhone != "" {
				err := s.Send(ctx, WhatsAppNotification{
					Type:     NotificationTypeTenant,
					TenantID: tenantID,
					Phone:    recipientPhone,
					Message:  message,
				})
				if err != nil {
					log.Printf("[ERROR] Failed to send tenant WA notification: %v", err)
				}
			}

		case model.NotificationChannelTelegram:
			// Telegram handled separately via existing telegram service
			log.Printf("[INFO] Telegram notification for event %s - handled by telegram service", eventType)
		}
	}

	return nil
}

// renderTemplate replaces {{variable}} placeholders with actual values
func (s *NotificationService) renderTemplate(content string, variables map[string]string) string {
	result := content
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
