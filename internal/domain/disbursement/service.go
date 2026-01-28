package disbursement_domain

import (
	"context"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type Service interface {
	GetAll(ctx context.Context) ([]model.Disbursement, error)
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id string, reason string) error
}

type service struct {
	repo outbound_port.DisbursementDatabasePort
}

func NewService(repo outbound_port.DisbursementDatabasePort) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAll(ctx context.Context) ([]model.Disbursement, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) Approve(ctx context.Context, id string) error {
	return s.repo.Approve(ctx, id)
}

func (s *service) Reject(ctx context.Context, id string, reason string) error {
	return s.repo.Reject(ctx, id, reason)
}
