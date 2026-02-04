package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type WhatsAppDatabasePort interface {
	GetByTenantID(ctx context.Context, tenantID string) (*model.WhatsAppSession, error)
	Save(ctx context.Context, session *model.WhatsAppSession) error
	UpdateStatus(ctx context.Context, sessionID, status string) error
	Delete(ctx context.Context, sessionID string) error
}
