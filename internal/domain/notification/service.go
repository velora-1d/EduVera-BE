package notification_domain

import (
	"context"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type Service interface {
	GetAll(ctx context.Context) ([]model.Notification, error)
	GetStats(ctx context.Context) (*model.NotificationStats, error)
}

type service struct {
	repo outbound_port.NotificationDatabasePort
}

func NewService(repo outbound_port.NotificationDatabasePort) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAll(ctx context.Context) ([]model.Notification, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetStats(ctx context.Context) (*model.NotificationStats, error) {
	return s.repo.GetStats(ctx)
}
