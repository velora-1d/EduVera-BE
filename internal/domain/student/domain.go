package student

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type StudentDomain interface {
	Create(ctx context.Context, tenantID string, input *model.StudentInput) (*model.Student, error)
	Update(ctx context.Context, tenantID, id string, input *model.StudentInput) (*model.Student, error)
	Delete(ctx context.Context, tenantID, id string) error
	FindByID(ctx context.Context, tenantID, id string) (*model.Student, error)
	List(ctx context.Context, filter model.StudentFilter) ([]model.Student, error)
	CountByTenant(ctx context.Context, tenantID string) (int, error)
}

type studentDomain struct {
	databasePort outbound_port.DatabasePort
}

func NewStudentDomain(databasePort outbound_port.DatabasePort) StudentDomain {
	return &studentDomain{
		databasePort: databasePort,
	}
}

func (d *studentDomain) Create(ctx context.Context, tenantID string, input *model.StudentInput) (*model.Student, error) {
	// Parse birth_date if provided
	var birthDate *time.Time
	if input.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil {
			birthDate = &parsed
		}
	}

	// Set default status
	status := input.Status
	if status == "" {
		status = model.StudentStatusActive
	}

	// Set default type
	studentType := input.Type
	if studentType == "" {
		studentType = model.StudentTypeSiswa
	}

	entryDate := time.Now()

	student := &model.Student{
		ID:               uuid.New().String(),
		TenantID:         tenantID,
		Name:             input.Name,
		Gender:           ptrString(input.Gender),
		BirthPlace:       ptrString(input.BirthPlace),
		BirthDate:        birthDate,
		NIS:              ptrString(input.NIS),
		NISN:             ptrString(input.NISN),
		NISPesantren:     ptrString(input.NISPesantren),
		Address:          ptrString(input.Address),
		Phone:            ptrString(input.Phone),
		FatherName:       ptrString(input.FatherName),
		FatherPhone:      ptrString(input.FatherPhone),
		FatherOccupation: ptrString(input.FatherOccupation),
		MotherName:       ptrString(input.MotherName),
		MotherPhone:      ptrString(input.MotherPhone),
		MotherOccupation: ptrString(input.MotherOccupation),
		GuardianName:     ptrString(input.GuardianName),
		GuardianPhone:    ptrString(input.GuardianPhone),
		GuardianRelation: ptrString(input.GuardianRelation),
		Type:             studentType,
		Jenjang:          ptrString(input.Jenjang),
		KelasID:          ptrString(input.KelasID),
		KamarID:          ptrString(input.KamarID),
		IsMukim:          input.IsMukim,
		TahunMasuk:       ptrInt(input.TahunMasuk),
		Status:           status,
		EntryDate:        &entryDate,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := d.databasePort.Student().Create(student)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to create student")
	}

	return student, nil
}

func (d *studentDomain) Update(ctx context.Context, tenantID, id string, input *model.StudentInput) (*model.Student, error) {
	// 1. Verify ownership securely by passing tenantID (which comes from JWT in Handler)
	// If tenantID doesn't match, FindByID will return error (Not Found)
	student, err := d.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find student")
	}

	// Update fields
	student.Name = input.Name
	if input.Gender != "" {
		student.Gender = &input.Gender
	}
	if input.BirthPlace != "" {
		student.BirthPlace = &input.BirthPlace
	}
	if input.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil {
			student.BirthDate = &parsed
		}
	}
	if input.NIS != "" {
		student.NIS = &input.NIS
	}
	if input.NISN != "" {
		student.NISN = &input.NISN
	}
	if input.NISPesantren != "" {
		student.NISPesantren = &input.NISPesantren
	}
	if input.Address != "" {
		student.Address = &input.Address
	}
	if input.Phone != "" {
		student.Phone = &input.Phone
	}
	if input.Type != "" {
		student.Type = input.Type
	}
	if input.Jenjang != "" {
		student.Jenjang = &input.Jenjang
	}
	if input.KelasID != "" {
		student.KelasID = &input.KelasID
	}
	if input.KamarID != "" {
		student.KamarID = &input.KamarID
	}
	student.IsMukim = input.IsMukim
	if input.Status != "" {
		student.Status = input.Status
	}
	student.UpdatedAt = time.Now()

	err = d.databasePort.Student().Update(student)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to update student")
	}

	return student, nil
}

func (d *studentDomain) Delete(ctx context.Context, tenantID, id string) error {
	return d.databasePort.Student().Delete(tenantID, id)
}

func (d *studentDomain) FindByID(ctx context.Context, tenantID, id string) (*model.Student, error) {
	return d.databasePort.Student().FindByID(tenantID, id)
}

func (d *studentDomain) List(ctx context.Context, filter model.StudentFilter) ([]model.Student, error) {
	return d.databasePort.Student().FindByFilter(filter)
}

func (d *studentDomain) CountByTenant(ctx context.Context, tenantID string) (int, error) {
	return d.databasePort.Student().CountByTenant(tenantID)
}

// Helper functions
func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
