package postgres_outbound_adapter

import (
	"database/sql"
	"fmt"
	"strings"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
)

const tableStudents = "students"

type studentAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewStudentAdapter(db outbound_port.DatabaseExecutor) outbound_port.StudentDatabasePort {
	return &studentAdapter{db: db}
}

func (a *studentAdapter) Create(student *model.Student) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tableStudents).Rows(student)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *studentAdapter) Update(student *model.Student) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableStudents).
		Set(student).
		Where(goqu.Ex{"id": student.ID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *studentAdapter) Delete(tenantID, id string) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Delete(tableStudents).Where(goqu.Ex{"id": id, "tenant_id": tenantID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *studentAdapter) FindByID(tenantID, id string) (*model.Student, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableStudents).Where(goqu.Ex{"id": id, "tenant_id": tenantID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	row := a.db.QueryRow(query)
	return a.scanStudent(row)
}

func (a *studentAdapter) FindByFilter(filter model.StudentFilter) ([]model.Student, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableStudents).Where(goqu.Ex{"tenant_id": filter.TenantID})

	if len(filter.IDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}
	if filter.Type != "" {
		if filter.Type == "both" {
			// For "both", we want students who are both siswa AND santri
			dataset = dataset.Where(goqu.Ex{"type": "both"})
		} else {
			// For siswa/santri, include those with that type OR "both"
			dataset = dataset.Where(goqu.Or(
				goqu.Ex{"type": filter.Type},
				goqu.Ex{"type": "both"},
			))
		}
	}
	if filter.Jenjang != "" {
		dataset = dataset.Where(goqu.Ex{"jenjang": filter.Jenjang})
	}
	if filter.KelasID != "" {
		dataset = dataset.Where(goqu.Ex{"kelas_id": filter.KelasID})
	}
	if filter.KamarID != "" {
		dataset = dataset.Where(goqu.Ex{"kamar_id": filter.KamarID})
	}
	if filter.IsMukim != nil {
		dataset = dataset.Where(goqu.Ex{"is_mukim": *filter.IsMukim})
	}
	if filter.Status != "" {
		dataset = dataset.Where(goqu.Ex{"status": filter.Status})
	}
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		dataset = dataset.Where(goqu.Or(
			goqu.L("LOWER(name)").Like(searchTerm),
			goqu.L("LOWER(nis)").Like(searchTerm),
			goqu.L("LOWER(nisn)").Like(searchTerm),
			goqu.L("LOWER(nis_pesantren)").Like(searchTerm),
		))
	}

	dataset = dataset.Order(goqu.I("name").Asc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		student, err := a.scanStudentRows(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, *student)
	}

	return students, nil
}

func (a *studentAdapter) CountByTenant(tenantID string) (int, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableStudents).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"tenant_id": tenantID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return 0, err
	}

	var count int
	err = a.db.QueryRow(query).Scan(&count)
	return count, err
}

func (a *studentAdapter) scanStudent(row *sql.Row) (*model.Student, error) {
	var s model.Student
	err := row.Scan(
		&s.ID, &s.TenantID,
		&s.Name, &s.Nickname, &s.Gender, &s.BirthPlace, &s.BirthDate, &s.PhotoURL,
		&s.NIS, &s.NISN, &s.NISPesantren,
		&s.Address, &s.Phone,
		&s.FatherName, &s.FatherPhone, &s.FatherOccupation,
		&s.MotherName, &s.MotherPhone, &s.MotherOccupation,
		&s.GuardianName, &s.GuardianPhone, &s.GuardianRelation,
		&s.Type, &s.Jenjang,
		&s.KelasID,
		&s.KamarID, &s.IsMukim, &s.TahunMasuk,
		&s.Status, &s.EntryDate, &s.ExitDate, &s.ExitReason,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan student: %w", err)
	}
	return &s, nil
}

func (a *studentAdapter) scanStudentRows(rows *sql.Rows) (*model.Student, error) {
	var s model.Student
	err := rows.Scan(
		&s.ID, &s.TenantID,
		&s.Name, &s.Nickname, &s.Gender, &s.BirthPlace, &s.BirthDate, &s.PhotoURL,
		&s.NIS, &s.NISN, &s.NISPesantren,
		&s.Address, &s.Phone,
		&s.FatherName, &s.FatherPhone, &s.FatherOccupation,
		&s.MotherName, &s.MotherPhone, &s.MotherOccupation,
		&s.GuardianName, &s.GuardianPhone, &s.GuardianRelation,
		&s.Type, &s.Jenjang,
		&s.KelasID,
		&s.KamarID, &s.IsMukim, &s.TahunMasuk,
		&s.Status, &s.EntryDate, &s.ExitDate, &s.ExitReason,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan student: %w", err)
	}
	return &s, nil
}
