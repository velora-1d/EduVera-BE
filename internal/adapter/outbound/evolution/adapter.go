package evolution

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type evolutionAdapter struct {
	baseURL   string
	globalKey string
}

func NewEvolutionAdapter() outbound_port.EvolutionApiPort {
	return &evolutionAdapter{
		baseURL:   os.Getenv("EVOLUTION_API_URL"),
		globalKey: os.Getenv("EVOLUTION_GLOBAL_KEY"),
	}
}

type createInstancePayload struct {
	InstanceName string `json:"instanceName"`
	Token        string `json:"token"`
	QRCode       bool   `json:"qrcode"`
}

type connectInstanceResponse struct {
	QRCode struct {
		Base64 string `json:"base64"`
	} `json:"qrcode"`
}

func (a *evolutionAdapter) CreateInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error) {
	url := fmt.Sprintf("%s/instance/create", a.baseURL)
	payload := createInstancePayload{
		InstanceName: instanceName,
		Token:        token,
		QRCode:       true,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", a.globalKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create instance: %s", string(bodyBytes))
	}

	return &model.WhatsAppSession{
		InstanceName: instanceName,
		APIKey:       token,
		Status:       model.WhatsAppStatusConnecting,
	}, nil
}

func (a *evolutionAdapter) FetchInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error) {
	url := fmt.Sprintf("%s/instance/connectionState/%s", a.baseURL, instanceName)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", a.globalKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch instance status")
	}

	var result struct {
		Instance struct {
			State string `json:"state"`
		} `json:"instance"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	var status string
	switch result.Instance.State {
	case "open":
		status = model.WhatsAppStatusConnected
	case "connecting":
		status = model.WhatsAppStatusConnecting
	default:
		status = model.WhatsAppStatusDisconnected
	}

	return &model.WhatsAppSession{
		InstanceName: instanceName,
		Status:       status,
	}, nil
}

func (a *evolutionAdapter) ConnectInstance(ctx context.Context, instanceName string, token string) (string, error) {
	// Usually CreateInstance returns QR if qrcode=true.
	// Or call /instance/connect/{instance}
	url := fmt.Sprintf("%s/instance/connect/%s", a.baseURL, instanceName)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", a.globalKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to connect instance")
	}

	var result connectInstanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.QRCode.Base64, nil
}

func (a *evolutionAdapter) LogoutInstance(ctx context.Context, instanceName string, token string) error {
	url := fmt.Sprintf("%s/instance/logout/%s", a.baseURL, instanceName)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("apikey", a.globalKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to logout instance")
	}

	return nil
}

func (a *evolutionAdapter) DeleteInstance(ctx context.Context, instanceName string, token string) error {
	url := fmt.Sprintf("%s/instance/delete/%s", a.baseURL, instanceName)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("apikey", a.globalKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete instance")
	}

	return nil
}

func (a *evolutionAdapter) SendMessage(ctx context.Context, instanceName string, token string, phone string, message string) error {
	url := fmt.Sprintf("%s/message/sendText/%s", a.baseURL, instanceName)

	payload := map[string]interface{}{
		"number": phone,
		"options": map[string]interface{}{
			"delay":    1200,
			"presence": "composing",
		},
		"textMessage": map[string]interface{}{
			"text": message,
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", a.globalKey) // Use global key for management, or token? Evolution v2 uses global key usually.

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s", string(bodyBytes))
	}

	return nil
}
