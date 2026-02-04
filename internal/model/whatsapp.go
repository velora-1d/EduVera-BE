package model

import (
	"time"
)

type WhatsAppSession struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	InstanceName string    `json:"instance_name" db:"instance_name"`
	APIKey       string    `json:"api_key" db:"api_key"`
	Status       string    `json:"status" db:"status"` // connecting, connected, disconnected
	QRCode       string    `json:"qr_code" db:"qr_code"`
	DeviceInfo   string    `json:"device_info" db:"device_info"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

const (
	WhatsAppStatusConnecting   = "connecting"
	WhatsAppStatusConnected    = "connected"
	WhatsAppStatusDisconnected = "disconnected"
)

type ConnectWhatsAppInput struct {
	TenantID string `json:"tenant_id" validate:"required"`
}

type WhatsAppStatusResponse struct {
	Status string `json:"status"`
	QRCode string `json:"qr_code,omitempty"`
}
