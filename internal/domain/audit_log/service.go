package audit_log

import (
	"context"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

// Service provides audit log operations
type Service interface {
	LogAction(ctx context.Context, input *model.AuditLogInput) error
	GetLogs(ctx context.Context, filter model.AuditLogFilter) ([]model.AuditLog, error)
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

type service struct {
	repo outbound_port.AuditLogDatabasePort
}

// NewService creates a new audit log service
func NewService(repo outbound_port.AuditLogDatabasePort) Service {
	return &service{repo: repo}
}

func (s *service) LogAction(ctx context.Context, input *model.AuditLogInput) error {
	log := &model.AuditLog{
		AdminID:     input.AdminID,
		AdminEmail:  input.AdminEmail,
		Action:      input.Action,
		TargetType:  input.TargetType,
		TargetID:    input.TargetID,
		OldValue:    input.OldValue,
		NewValue:    input.NewValue,
		IPAddress:   input.IPAddress,
		UserAgent:   input.UserAgent,
		Description: input.Description,
	}

	return s.repo.Create(log)
}

func (s *service) GetLogs(ctx context.Context, filter model.AuditLogFilter) ([]model.AuditLog, error) {
	return s.repo.FindByFilter(filter)
}

func (s *service) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetStats()
}
