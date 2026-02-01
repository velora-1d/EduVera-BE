package spp_domain

import (
	"context"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type Service interface {
	ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error)
	Create(ctx context.Context, spp *model.SPPTransaction) error
	RecordPayment(ctx context.Context, id string, paymentMethod string) error
	GetStats(ctx context.Context, tenantID string) (*model.SPPStats, error)
	ListAll(ctx context.Context) ([]model.SPPTransaction, error)
	// Manual payment methods
	Update(ctx context.Context, id, studentName string, amount int64, description, dueDate, period string) error
	Delete(ctx context.Context, id string) error
	UploadProof(ctx context.Context, id string, proofURL string) error
	ConfirmPayment(ctx context.Context, id string, confirmedBy string, paymentMethod string) error
	ListOverdue(ctx context.Context, tenantID string) ([]model.SPPTransaction, error)
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

// Update modifies an SPP transaction
func (s *service) Update(ctx context.Context, id, studentName string, amount int64, description, dueDateStr, period string) error {
	spp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	spp.StudentName = studentName
	spp.Amount = amount
	spp.Description = description
	spp.Period = period

	// Parse due date if provided
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			spp.DueDate = &dueDate
		}
	}

	return s.repo.Update(ctx, spp)
}

// Delete removes an SPP transaction
func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// UploadProof saves the payment proof URL
func (s *service) UploadProof(ctx context.Context, id string, proofURL string) error {
	return s.repo.UploadProof(ctx, id, proofURL)
}

// ConfirmPayment marks payment as confirmed by admin
func (s *service) ConfirmPayment(ctx context.Context, id string, confirmedBy string, paymentMethod string) error {
	// First update payment method if provided
	if paymentMethod != "" {
		if err := s.repo.UpdateStatus(ctx, id, model.SPPStatusPaid, paymentMethod); err != nil {
			return err
		}
	}
	// Then mark as confirmed
	return s.repo.ConfirmPayment(ctx, id, confirmedBy)
}

// ListOverdue returns overdue pending payments
func (s *service) ListOverdue(ctx context.Context, tenantID string) ([]model.SPPTransaction, error) {
	return s.repo.ListOverdue(ctx, tenantID)
}
