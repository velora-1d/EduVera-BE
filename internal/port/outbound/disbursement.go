package outbound_port

import (
	"context"
	"eduvera/internal/model"
)

type DisbursementDatabasePort interface {
	GetAll(ctx context.Context) ([]model.Disbursement, error)
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id string, reason string) error
}
