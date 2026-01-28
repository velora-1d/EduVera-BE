package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TelegramNotifier sends notifications to Telegram
type TelegramNotifier struct {
	botToken string
	chatID   string
	client   *http.Client
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier() *TelegramNotifier {
	return &TelegramNotifier{
		botToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegistrationData contains data for new registration notification
type RegistrationData struct {
	InstitutionName string
	PlanType        string
	AdminName       string
	Email           string
	WhatsApp        string
	Subdomain       string
	Address         string
}

// SendNewRegistration sends a notification for new registration
func (t *TelegramNotifier) SendNewRegistration(data RegistrationData) error {
	planEmoji := map[string]string{
		"sekolah":   "🏫",
		"pesantren": "🕌",
		"hybrid":    "🏛️",
	}

	emoji := planEmoji[data.PlanType]
	if emoji == "" {
		emoji = "📌"
	}

	message := fmt.Sprintf(`🎉 *PENDAFTARAN BARU!*

%s *Paket:* %s

📌 *Lembaga:* %s
📍 *Alamat:* %s
🌐 *Subdomain:* %s.eduvera.ve-lora.my.id

👤 *Admin:* %s
📧 *Email:* %s
📱 *WhatsApp:* %s

🕐 *Waktu:* %s`,
		emoji,
		data.PlanType,
		data.InstitutionName,
		data.Address,
		data.Subdomain,
		data.AdminName,
		data.Email,
		data.WhatsApp,
		time.Now().Format("02 Jan 2006, 15:04 WIB"),
	)

	return t.sendMessage(message)
}

// SendPaymentSuccess sends a notification for successful payment
func (t *TelegramNotifier) SendPaymentSuccess(institutionName, planType string, amount int64) error {
	message := fmt.Sprintf(`💰 *PEMBAYARAN BERHASIL!*

🏫 *Lembaga:* %s
📦 *Paket:* %s
💵 *Jumlah:* Rp %d

✅ Tenant berhasil diaktifkan!`,
		institutionName,
		planType,
		amount,
	)

	return t.sendMessage(message)
}

// SendTestMessage sends a test message
func (t *TelegramNotifier) SendTestMessage() error {
	message := `🚀 *EduVera Bot Connected!*

Bot notifikasi EduVera berhasil terhubung.
Anda akan menerima notif di sini setiap ada:
- 📝 Pendaftaran baru
- 💰 Pembayaran berhasil
- ⚠️ Alert penting

_Powered by EduVera SaaS_`

	return t.sendMessage(message)
}

// sendMessage sends a message via Telegram Bot API
func (t *TelegramNotifier) sendMessage(text string) error {
	if t.botToken == "" || t.chatID == "" {
		return fmt.Errorf("telegram config not set")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	payload := map[string]interface{}{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := t.client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %d", resp.StatusCode)
	}

	return nil
}
