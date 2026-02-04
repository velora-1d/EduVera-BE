package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// AlertService sends alerts to owner via Telegram
type AlertService struct {
	botToken string
	chatID   string
}

func NewAlertService() *AlertService {
	return &AlertService{
		botToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}
}

// SendTelegramAlert sends alert message to owner
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

// AlertPaymentFailed sends alert when payment fails
func (s *AlertService) AlertPaymentFailed(tenantName, amount, orderId, reason string) error {
	message := fmt.Sprintf(
		"🚨 <b>Payment Failed</b>\n\n"+
			"<b>Tenant:</b> %s\n"+
			"<b>Amount:</b> Rp %s\n"+
			"<b>Order ID:</b> %s\n"+
			"<b>Reason:</b> %s\n\n"+
			"Please check dashboard for details.",
		tenantName, amount, orderId, reason,
	)
	return s.SendTelegramAlert(message)
}

// AlertSubscriptionExpiring sends alert for expiring subscriptions
func (s *AlertService) AlertSubscriptionExpiring(tenantName, expiryDate string, daysLeft int) error {
	message := fmt.Sprintf(
		"⚠️ <b>Subscription Expiring</b>\n\n"+
			"<b>Tenant:</b> %s\n"+
			"<b>Expires:</b> %s\n"+
			"<b>Days Left:</b> %d\n\n"+
			"Reminder sent to tenant via WhatsApp.",
		tenantName, expiryDate, daysLeft,
	)
	return s.SendTelegramAlert(message)
}

// AlertNewTenantRegistered sends alert for new registrations
func (s *AlertService) AlertNewTenantRegistered(tenantName, ownerEmail, planType string) error {
	message := fmt.Sprintf(
		"🎉 <b>New Tenant Registered</b>\n\n"+
			"<b>Name:</b> %s\n"+
			"<b>Owner:</b> %s\n"+
			"<b>Plan:</b> %s\n\n"+
			"Welcome to EduVera!",
		tenantName, ownerEmail, planType,
	)
	return s.SendTelegramAlert(message)
}

// AlertSystemError sends alert for critical system errors
func (s *AlertService) AlertSystemError(component, errorMsg string) error {
	message := fmt.Sprintf(
		"🔴 <b>System Error</b>\n\n"+
			"<b>Component:</b> %s\n"+
			"<b>Error:</b> %s\n\n"+
			"Immediate attention required!",
		component, errorMsg,
	)
	return s.SendTelegramAlert(message)
}
