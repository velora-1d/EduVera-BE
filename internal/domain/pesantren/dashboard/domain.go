package dashboard

import (
	"context"

	"prabogo/internal/model"
)

type Service interface {
	GetStats(ctx context.Context, tenantID string) (*model.DashboardStats, error)
}
