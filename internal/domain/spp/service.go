package spp_domain

import (
	"context"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type Service interface {
	ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error)
	Create(ctx context.Context, spp *model.SPPTransaction) error
	RecordPayment(ctx context.Context, id string, paymentMethod string) error
	GetStats(ctx context.Context, tenantID string) (*model.SPPStats, error)
	ListAll(ctx context.Context) ([]model.SPPTransaction, error)
}

type service struct {
	repo outbound_port.SPPDatabasePort
}

func NewService(repo outbound_port.SPPDatabasePort) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *service) Create(ctx context.Context, spp *model.SPPTransaction) error {
	spp.Status = model.SPPStatusPending
	return s.repo.Create(ctx, spp)
}

func (s *service) RecordPayment(ctx context.Context, id string, paymentMethod string) error {
	return s.repo.UpdateStatus(ctx, id, model.SPPStatusPaid, paymentMethod)
}

func (s *service) GetStats(ctx context.Context, tenantID string) (*model.SPPStats, error) {
	return s.repo.GetStatsByTenant(ctx, tenantID)
}

func (s *service) ListAll(ctx context.Context) ([]model.SPPTransaction, error) {
	return s.repo.ListAll(ctx)
}
