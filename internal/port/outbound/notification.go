package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type NotificationDatabasePort interface {
	GetAll(ctx context.Context) ([]model.Notification, error)
	GetStats(ctx context.Context) (*model.NotificationStats, error)
	Create(ctx context.Context, notification *model.Notification) error
}
