package whatsapp_outbound_adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	outbound_port "eduvera/internal/port/outbound"
)

type fonnteAdapter struct {
	apiKey string
}

func NewFonnteAdapter() outbound_port.WhatsAppMessagePort {
	return &fonnteAdapter{
		apiKey: os.Getenv("FONNTE_TOKEN"),
	}
}

func (a *fonnteAdapter) Send(target, message string) error {
	url := "https://api.fonnte.com/send"

	payload := map[string]string{
		"target":  target,
		"message": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fonnte api returned status: %d", resp.StatusCode)
	}

	return nil
}
