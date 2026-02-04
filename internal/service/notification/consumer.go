package notification

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"prabogo/utils/rabbitmq"
)

const (
	WhatsAppExchange = "whatsapp_notifications"
	WhatsAppQueue    = "whatsapp_queue"
)

// StartWhatsAppConsumer starts the message consumer for WhatsApp notifications
func (s *NotificationService) StartWhatsAppConsumer(ctx context.Context) error {
	log.Println("[INFO] Starting WhatsApp notification consumer...")

	cfg := rabbitmq.SubscriberConfig{
		Exchange:     WhatsAppExchange,
		ExchangeKind: rabbitmq.KindFanOut,
		Queue:        WhatsAppQueue,
		RouteKey:     "",
		Callback: func(msg []byte) bool {
			return s.handleMessage(ctx, msg)
		},
	}

	return rabbitmq.SubscriberWithConfig(cfg)
}

// handleMessage processes a single notification message
func (s *NotificationService) handleMessage(ctx context.Context, msg []byte) bool {
	var notif WhatsAppNotification
	if err := json.Unmarshal(msg, &notif); err != nil {
		log.Printf("[ERROR] Failed to unmarshal notification: %v", err)
		return true // Ack to prevent infinite requeue of invalid messages
	}

	// Apply rate limiting (simple delay)
	time.Sleep(100 * time.Millisecond)

	err := s.Send(ctx, notif)
	if err != nil {
		log.Printf("[ERROR] Failed to send notification to %s: %v", notif.Phone, err)

		// Retry logic
		if notif.RetryLeft > 0 {
			notif.RetryLeft--
			_ = PublishNotification(notif)
			return true // Ack current message, requeued with decremented retry
		}

		return true // Max retries reached, ack and drop
	}

	log.Printf("[INFO] Notification sent to %s via %s", notif.Phone, notif.Type)
	return true
}

// PublishNotification publishes a notification to the queue
func PublishNotification(notif WhatsAppNotification) error {
	ctx := context.Background()

	// Set default retry count if not specified
	if notif.RetryLeft == 0 {
		notif.RetryLeft = 3
	}

	return rabbitmq.Publish(ctx, WhatsAppExchange, rabbitmq.KindFanOut, "", notif)
}

// PublishSystemNotification helper for system-level notifications (via Fonnte)
func PublishSystemNotification(phone, message string) error {
	return PublishNotification(WhatsAppNotification{
		Type:    NotificationTypeSystem,
		Phone:   phone,
		Message: message,
	})
}

// PublishTenantNotification helper for tenant-level notifications (via Evolution API)
func PublishTenantNotification(tenantID, phone, message string) error {
	return PublishNotification(WhatsAppNotification{
		Type:     NotificationTypeTenant,
		TenantID: tenantID,
		Phone:    phone,
		Message:  message,
	})
}
