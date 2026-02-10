package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type SPPDatabasePort interface {
	Create(ctx context.Context, spp *model.SPPTransaction) error
	Update(ctx context.Context, spp *model.SPPTransaction) error
	// Delete(ctx context.Context, tenantID, id string) error
	Delete(ctx context.Context, tenantID, id string) error
	ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error)
	FindByID(ctx context.Context, tenantID, id string) (*model.SPPTransaction, error)
	UpdateStatus(ctx context.Context, tenantID, id string, status model.SPPStatus, paymentMethod string) error
	GetStatsByTenant(ctx context.Context, tenantID string) (*model.SPPStats, error)
	ListAll(ctx context.Context) ([]model.SPPTransaction, error)
	// Manual payment confirmation
	UploadProof(ctx context.Context, tenantID, id string, proofURL string) error
	ConfirmPayment(ctx context.Context, tenantID, id string, confirmedBy string) error
	ListByPeriod(ctx context.Context, tenantID string, period string) ([]model.SPPTransaction, error)
	ListOverdue(ctx context.Context, tenantID string) ([]model.SPPTransaction, error)
}
