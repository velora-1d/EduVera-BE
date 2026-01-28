package dashboard

import (
	"context"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type service struct {
	db outbound_port.PesantrenDashboardPort
}

func NewService(db outbound_port.PesantrenDashboardPort) Service {
	return &service{
		db: db,
	}
}

func (s *service) GetStats(ctx context.Context, tenantID string) (*model.DashboardStats, error) {
	return s.db.GetDashboardStats(ctx, tenantID)
}
