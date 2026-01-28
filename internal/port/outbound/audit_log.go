package outbound_port

import "eduvera/internal/model"

// AuditLogDatabasePort defines the interface for audit log operations
type AuditLogDatabasePort interface {
	Create(log *model.AuditLog) error
	FindByFilter(filter model.AuditLogFilter) ([]model.AuditLog, error)
	GetStats() (map[string]interface{}, error)
}
