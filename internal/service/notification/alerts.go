package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	outbound_port "prabogo/internal/port/outbound"
)

// AlertService sends alerts to owner via Telegram and WhatsApp
type AlertService struct {
	botToken     string
	chatID       string
	notifService *NotificationService
}

func NewAlertService() *AlertService {
	return &AlertService{
		botToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}
}

// NewAlertServiceWithWA creates AlertService with WhatsApp notification support
func NewAlertServiceWithWA(
	fonntePort outbound_port.WhatsAppMessagePort,
	evolutionPort outbound_port.WhatsAppClientPort,
	dbPort outbound_port.DatabasePort,
) *AlertService {
	return &AlertService{
		botToken:     os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatID:       os.Getenv("TELEGRAM_CHAT_ID"),
		notifService: NewNotificationService(fonntePort, evolutionPort, dbPort),
	}
}

// SendTelegramAlert sends alert message to owner via Telegram
func (s *AlertService) SendTelegramAlert(message string) error {
	if s.botToken == "" || s.chatID == "" {
		log.Println("[WARN] Telegram not configured, skipping alert")
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	payload := map[string]interface{}{
		"chat_id":    s.chatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status: %d", resp.StatusCode)
	}

	return nil
}

// sendOwnerWA sends notification to owner's WhatsApp
func (s *AlertService) sendOwnerWA(ctx context.Context, message string) error {
	if s.notifService == nil {
		log.Println("[WARN] WhatsApp notification service not initialized, skipping WA alert")
		return nil
	}

	ownerPhone := os.Getenv("OWNER_PHONE")
	if ownerPhone == "" {
		ownerPhone = "6285117776596" // Default owner phone
	}

	return s.notifService.Send(ctx, WhatsAppNotification{
		Type:    NotificationTypeOwner,
		Phone:   ownerPhone,
		Message: message,
	})
}

// sendMultiChannel sends message to both Telegram and WhatsApp
func (s *AlertService) sendMultiChannel(ctx context.Context, telegramMsg, waMsg string) error {
	// Send Telegram (non-blocking)
	go func() {
		if err := s.SendTelegramAlert(telegramMsg); err != nil {
			log.Printf("[ERROR] Failed to send Telegram alert: %v", err)
		}
	}()

	// Send WhatsApp
	if err := s.sendOwnerWA(ctx, waMsg); err != nil {
		log.Printf("[ERROR] Failed to send Owner WA alert: %v", err)
		return err
	}

	return nil
}

// AlertPaymentFailed sends alert when payment fails
func (s *AlertService) AlertPaymentFailed(tenantName, amount, orderId, reason string) error {
	telegramMsg := fmt.Sprintf(
		"🚨 <b>Payment Failed</b>\n\n"+
			"<b>Tenant:</b> %s\n"+
			"<b>Amount:</b> Rp %s\n"+
			"<b>Order ID:</b> %s\n"+
			"<b>Reason:</b> %s\n\n"+
			"Please check dashboard for details.",
		tenantName, amount, orderId, reason,
	)

	waMsg := fmt.Sprintf(
		"🚨 *Payment Failed*\n\n"+
			"Tenant: %s\n"+
			"Amount: Rp %s\n"+
			"Order ID: %s\n"+
			"Reason: %s",
		tenantName, amount, orderId, reason,
	)

	return s.sendMultiChannel(context.Background(), telegramMsg, waMsg)
}

// AlertSubscriptionExpiring sends alert for expiring subscriptions
func (s *AlertService) AlertSubscriptionExpiring(tenantName, expiryDate string, daysLeft int) error {
	telegramMsg := fmt.Sprintf(
		"⚠️ <b>Subscription Expiring</b>\n\n"+
			"<b>Tenant:</b> %s\n"+
			"<b>Expires:</b> %s\n"+
			"<b>Days Left:</b> %d\n\n"+
			"Reminder sent to tenant via WhatsApp.",
		tenantName, expiryDate, daysLeft,
	)

	waMsg := fmt.Sprintf(
		"⚠️ *Subscription Expiring*\n\n"+
			"Tenant: %s\n"+
			"Expires: %s\n"+
			"Days Left: %d",
		tenantName, expiryDate, daysLeft,
	)

	return s.sendMultiChannel(context.Background(), telegramMsg, waMsg)
}

// AlertNewTenantRegistered sends alert for new registrations
func (s *AlertService) AlertNewTenantRegistered(tenantName, ownerEmail, planType string) error {
	telegramMsg := fmt.Sprintf(
		"🎉 <b>New Tenant Registered</b>\n\n"+
			"<b>Name:</b> %s\n"+
			"<b>Owner:</b> %s\n"+
			"<b>Plan:</b> %s\n\n"+
			"Welcome to EduVera!",
		tenantName, ownerEmail, planType,
	)

	waMsg := fmt.Sprintf(
		"🎉 *Pendaftaran Baru!*\n\n"+
			"Nama: %s\n"+
			"Owner: %s\n"+
			"Paket: %s\n\n"+
			"Selamat datang di EduVera!",
		tenantName, ownerEmail, planType,
	)

	return s.sendMultiChannel(context.Background(), telegramMsg, waMsg)
}

// AlertSystemError sends alert for critical system errors
func (s *AlertService) AlertSystemError(component, errorMsg string) error {
	telegramMsg := fmt.Sprintf(
		"🔴 <b>System Error</b>\n\n"+
			"<b>Component:</b> %s\n"+
			"<b>Error:</b> %s\n\n"+
			"Immediate attention required!",
		component, errorMsg,
	)

	waMsg := fmt.Sprintf(
		"🔴 *System Error*\n\n"+
			"Component: %s\n"+
			"Error: %s\n\n"+
			"Immediate attention required!",
		component, errorMsg,
	)

	return s.sendMultiChannel(context.Background(), telegramMsg, waMsg)
}

// AlertUpgrade sends alert when tenant upgrades subscription
func (s *AlertService) AlertUpgrade(tenantName, fromTier, toTier, amount string) error {
	telegramMsg := fmt.Sprintf(
		"💎 <b>Upgrade Subscription</b>\n\n"+
			"<b>Tenant:</b> %s\n"+
			"<b>From:</b> %s\n"+
			"<b>To:</b> %s\n"+
			"<b>Amount:</b> Rp %s\n\n"+
			"Thank you for upgrading!",
		tenantName, fromTier, toTier, amount,
	)

	waMsg := fmt.Sprintf(
		"💎 *Upgrade Paket!*\n\n"+
			"Tenant: %s\n"+
			"Dari: %s\n"+
			"Ke: %s\n"+
			"Nominal: Rp %s\n\n"+
			"Terima kasih telah upgrade!",
		tenantName, fromTier, toTier, amount,
	)

	return s.sendMultiChannel(context.Background(), telegramMsg, waMsg)
}

// AlertPaymentSuccess sends alert when payment succeeds
func (s *AlertService) AlertPaymentSuccess(tenantName, amount, orderId, paymentType string) error {
	telegramMsg := fmt.Sprintf(
		"💰 <b>Payment Received</b>\n\n"+
			"<b>Tenant:</b> %s\n"+
			"<b>Amount:</b> Rp %s\n"+
			"<b>Order ID:</b> %s\n"+
			"<b>Type:</b> %s",
		tenantName, amount, orderId, paymentType,
	)

	waMsg := fmt.Sprintf(
		"💰 *Pembayaran Diterima!*\n\n"+
			"Tenant: %s\n"+
			"Nominal: Rp %s\n"+
			"Order ID: %s\n"+
			"Tipe: %s",
		tenantName, amount, orderId, paymentType,
	)

	return s.sendMultiChannel(context.Background(), telegramMsg, waMsg)
}
