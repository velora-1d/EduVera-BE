package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type EvolutionApiPort interface {
	CreateInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error)
	FetchInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error)
	ConnectInstance(ctx context.Context, instanceName string, token string) (string, error) // Returns Base64 QR
	LogoutInstance(ctx context.Context, instanceName string, token string) error
	DeleteInstance(ctx context.Context, instanceName string, token string) error
	SendMessage(ctx context.Context, instanceName string, token string, phone string, message string) error
}
