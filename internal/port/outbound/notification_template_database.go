package outbound_port

import (
	"context"

	"prabogo/internal/model"
)

// NotificationTemplateDatabasePort defines database operations for notification templates
type NotificationTemplateDatabasePort interface {
	GetAll(ctx context.Context) ([]model.NotificationTemplate, error)
	GetByID(ctx context.Context, id string) (*model.NotificationTemplate, error)
	GetByEventAndChannel(ctx context.Context, eventType, channel string) (*model.NotificationTemplate, error)
	GetActiveByEvent(ctx context.Context, eventType string) ([]model.NotificationTemplate, error)
	Save(ctx context.Context, template *model.NotificationTemplate) error
	Update(ctx context.Context, template *model.NotificationTemplate) error
	Delete(ctx context.Context, id string) error
}
