package outbound_port

import "prabogo/internal/model"

// StudentDatabasePort defines the interface for student database operations
type StudentDatabasePort interface {
	Create(student *model.Student) error
	Update(student *model.Student) error
	Delete(id string) error
	FindByID(id string) (*model.Student, error)
	FindByFilter(filter model.StudentFilter) ([]model.Student, error)
	CountByTenant(tenantID string) (int, error)
}
