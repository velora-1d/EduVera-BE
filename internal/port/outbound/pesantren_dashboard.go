package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type PesantrenDashboardPort interface {
	GetDashboardStats(ctx context.Context, tenantID string) (*model.DashboardStats, error)
}
