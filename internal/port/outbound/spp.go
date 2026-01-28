package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type SPPDatabasePort interface {
	Create(ctx context.Context, spp *model.SPPTransaction) error
	ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error)
	FindByID(ctx context.Context, id string) (*model.SPPTransaction, error)
	UpdateStatus(ctx context.Context, id string, status model.SPPStatus, paymentMethod string) error
	GetStatsByTenant(ctx context.Context, tenantID string) (*model.SPPStats, error)
	ListAll(ctx context.Context) ([]model.SPPTransaction, error)
}
