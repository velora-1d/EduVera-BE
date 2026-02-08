package whatsback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type whatsbackAdapter struct {
	baseURL string
}

func NewWhatsbackAdapter() outbound_port.WhatsAppClientPort {
	return &whatsbackAdapter{
		baseURL: os.Getenv("EVOLUTION_API_URL"), // Reset to point to Whatsback
	}
}

type sessionResponse struct {
	Status bool `json:"status"`
	Data   struct {
		Status       string `json:"status"`
		QRCode       string `json:"qrcode"` // Base64
		PhoneNumber  string `json:"phone_number"`
		InstanceName string `json:"instance_name"`
	} `json:"data"`
}

func (a *whatsbackAdapter) CreateInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error) {
	// Whatsback is single session and auto-started.
	// We just check if it's reachable.
	_, err := a.fetchSessionStatus()
	if err != nil {
		return nil, err
	}

	return &model.WhatsAppSession{
		InstanceName: instanceName,
		APIKey:       token,
		Status:       model.WhatsAppStatusConnecting,
	}, nil
}

func (a *whatsbackAdapter) FetchInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error) {
	data, err := a.fetchSessionStatus()
	if err != nil {
		return nil, err
	}

	var status string
	switch data.Status {
	case "connected", "open":
		status = model.WhatsAppStatusConnected
	case "connecting":
		status = model.WhatsAppStatusConnecting
	default:
		status = model.WhatsAppStatusDisconnected
	}

	return &model.WhatsAppSession{
		InstanceName: instanceName,
		Status:       status,
		PhoneNumber:  data.PhoneNumber,
	}, nil
}

func (a *whatsbackAdapter) ConnectInstance(ctx context.Context, instanceName string, token string) (string, error) {
	// Returns QR Code
	data, err := a.fetchSessionStatus()
	if err != nil {
		return "", err
	}
	return data.QRCode, nil // Base64 or null
}

func (a *whatsbackAdapter) LogoutInstance(ctx context.Context, instanceName string, token string) error {
	url := fmt.Sprintf("%s/session/logout", a.baseURL)
	req, _ := http.NewRequest("POST", url, nil)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to logout")
	}
	return nil
}

func (a *whatsbackAdapter) DeleteInstance(ctx context.Context, instanceName string, token string) error {
	// Treat delete as logout for single session
	return a.LogoutInstance(ctx, instanceName, token)
}

func (a *whatsbackAdapter) SendMessage(ctx context.Context, instanceName string, token string, phone string, message string) error {
	url := fmt.Sprintf("%s/message/send-message", a.baseURL)

	// Ensure phone has + prefix for Whatsback
	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}
	// Remove @s.whatsapp.net if present (Evolution usually adds it, but input here is usually raw phone)
	phone = strings.ReplaceAll(phone, "@s.whatsapp.net", "")

	payload := map[string]string{
		"number":  phone,
		"message": message,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s", string(bodyBytes))
	}

	return nil
}

// Helper to get status from our custom endpoint
func (a *whatsbackAdapter) fetchSessionStatus() (*struct {
	Status       string `json:"status"`
	QRCode       string `json:"qrcode"`
	PhoneNumber  string `json:"phone_number"`
	InstanceName string `json:"instance_name"`
}, error) {
	url := fmt.Sprintf("%s/session/status", a.baseURL)
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch status")
	}

	var result sessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("api error")
	}

	return &result.Data, nil
}
