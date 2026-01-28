package outbound_port

import (
	"context"

	"prabogo/internal/model"
)

// ERaporDatabasePort defines the interface for E-Rapor database operations
type ERaporDatabasePort interface {
	// Subject operations
	CreateSubject(ctx context.Context, input *model.SubjectInput) (*model.Subject, error)
	UpdateSubject(ctx context.Context, id string, input *model.SubjectInput) (*model.Subject, error)
	GetSubjectByID(ctx context.Context, id string) (*model.Subject, error)
	GetSubjectsByTenant(ctx context.Context, tenantID string) ([]model.Subject, error)
	DeleteSubject(ctx context.Context, id string) error

	// Grade operations
	SaveGrade(ctx context.Context, input *model.StudentGradeInput) (*model.StudentGrade, error)
	BatchSaveGrades(ctx context.Context, input *model.BatchGradeInput) ([]model.StudentGrade, error)
	GetGradesByStudent(ctx context.Context, studentID, semesterID string) ([]model.StudentGrade, error)
	GetGradesBySubject(ctx context.Context, subjectID, semesterID string) ([]model.StudentGrade, error)
	GetStudentRapor(ctx context.Context, studentID, semesterID string) (*model.RaporData, error)

	// Stats
	GetGradeStats(ctx context.Context, tenantID, semesterID string) (map[string]interface{}, error)

	// Snapshot Operations (Rapor)
	GetOrCreateRaporPeriode(tenantID, name string) (*model.RaporPeriode, error)
	CreateRapor(m *model.Rapor) error
	CreateRaporNilai(m *model.RaporNilai) error
}
