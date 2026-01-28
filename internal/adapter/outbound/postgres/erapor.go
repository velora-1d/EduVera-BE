package postgres_outbound_adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type eraporAdapter struct {
	db *goqu.Database
}

func NewERaporAdapter(sqlDB *sql.DB) *eraporAdapter {
	return &eraporAdapter{db: goqu.New("postgres", sqlDB)}
}

// ==========================================
// SUBJECT OPERATIONS
// ==========================================

func (a *eraporAdapter) CreateSubject(ctx context.Context, input *model.SubjectInput) (*model.Subject, error) {
	id := uuid.New().String()
	now := time.Now()

	// Marshal grading config to JSON
	configJSON, err := json.Marshal(input.GradingConfig)
	if err != nil {
		return nil, err
	}

	_, err = a.db.Insert("subjects").Rows(
		goqu.Record{
			"id":             id,
			"tenant_id":      input.TenantID,
			"name":           input.Name,
			"code":           input.Code,
			"type":           input.Type,
			"grading_config": configJSON,
			"is_active":      true,
			"created_at":     now,
			"updated_at":     now,
		},
	).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Subject{
		ID:            id,
		TenantID:      input.TenantID,
		Name:          input.Name,
		Code:          input.Code,
		Type:          input.Type,
		GradingConfig: input.GradingConfig,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (a *eraporAdapter) UpdateSubject(ctx context.Context, id string, input *model.SubjectInput) (*model.Subject, error) {
	now := time.Now()

	configJSON, err := json.Marshal(input.GradingConfig)
	if err != nil {
		return nil, err
	}

	_, err = a.db.Update("subjects").Set(
		goqu.Record{
			"name":           input.Name,
			"code":           input.Code,
			"type":           input.Type,
			"grading_config": configJSON,
			"updated_at":     now,
		},
	).Where(goqu.C("id").Eq(id)).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return a.GetSubjectByID(ctx, id)
}

func (a *eraporAdapter) GetSubjectByID(ctx context.Context, id string) (*model.Subject, error) {
	var subject struct {
		ID            string    `db:"id"`
		TenantID      string    `db:"tenant_id"`
		Name          string    `db:"name"`
		Code          string    `db:"code"`
		Type          string    `db:"type"`
		GradingConfig []byte    `db:"grading_config"`
		IsActive      bool      `db:"is_active"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}

	found, err := a.db.From("subjects").Where(goqu.C("id").Eq(id)).ScanStructContext(ctx, &subject)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	var config model.GradingConfig
	if len(subject.GradingConfig) > 0 {
		json.Unmarshal(subject.GradingConfig, &config)
	}

	return &model.Subject{
		ID:            subject.ID,
		TenantID:      subject.TenantID,
		Name:          subject.Name,
		Code:          subject.Code,
		Type:          subject.Type,
		GradingConfig: config,
		IsActive:      subject.IsActive,
		CreatedAt:     subject.CreatedAt,
		UpdatedAt:     subject.UpdatedAt,
	}, nil
}

func (a *eraporAdapter) GetSubjectsByTenant(ctx context.Context, tenantID string) ([]model.Subject, error) {
	var rows []struct {
		ID            string    `db:"id"`
		TenantID      string    `db:"tenant_id"`
		Name          string    `db:"name"`
		Code          string    `db:"code"`
		Type          string    `db:"type"`
		GradingConfig []byte    `db:"grading_config"`
		IsActive      bool      `db:"is_active"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}

	err := a.db.From("subjects").
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Where(goqu.C("is_active").Eq(true)).
		Order(goqu.C("name").Asc()).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, err
	}

	subjects := make([]model.Subject, len(rows))
	for i, row := range rows {
		var config model.GradingConfig
		if len(row.GradingConfig) > 0 {
			json.Unmarshal(row.GradingConfig, &config)
		}
		subjects[i] = model.Subject{
			ID:            row.ID,
			TenantID:      row.TenantID,
			Name:          row.Name,
			Code:          row.Code,
			Type:          row.Type,
			GradingConfig: config,
			IsActive:      row.IsActive,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
	}

	return subjects, nil
}

func (a *eraporAdapter) DeleteSubject(ctx context.Context, id string) error {
	// Soft delete
	_, err := a.db.Update("subjects").Set(
		goqu.Record{"is_active": false, "updated_at": time.Now()},
	).Where(goqu.C("id").Eq(id)).Executor().ExecContext(ctx)
	return err
}

// ==========================================
// GRADE OPERATIONS
// ==========================================

func (a *eraporAdapter) SaveGrade(ctx context.Context, input *model.StudentGradeInput) (*model.StudentGrade, error) {
	id := uuid.New().String()
	now := time.Now()

	componentScoresJSON, _ := json.Marshal(input.ComponentScores)

	_, err := a.db.Insert("student_grades").Rows(
		goqu.Record{
			"id":               id,
			"tenant_id":        input.TenantID,
			"student_id":       input.StudentID,
			"subject_id":       input.SubjectID,
			"semester_id":      input.SemesterID,
			"score_numeric":    input.ScoreNumeric,
			"score_predicate":  input.ScorePredicate,
			"description_high": input.DescriptionHigh,
			"description_low":  input.DescriptionLow,
			"component_scores": componentScoresJSON,
			"created_at":       now,
			"updated_at":       now,
		},
	).OnConflict(
		goqu.DoUpdate("student_id, subject_id, semester_id",
			goqu.Record{
				"score_numeric":    input.ScoreNumeric,
				"score_predicate":  input.ScorePredicate,
				"description_high": input.DescriptionHigh,
				"description_low":  input.DescriptionLow,
				"component_scores": componentScoresJSON,
				"updated_at":       now,
			}),
	).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return &model.StudentGrade{
		ID:              id,
		TenantID:        input.TenantID,
		StudentID:       input.StudentID,
		SubjectID:       input.SubjectID,
		SemesterID:      input.SemesterID,
		ScoreNumeric:    input.ScoreNumeric,
		ScorePredicate:  input.ScorePredicate,
		DescriptionHigh: input.DescriptionHigh,
		DescriptionLow:  input.DescriptionLow,
		ComponentScores: input.ComponentScores,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (a *eraporAdapter) BatchSaveGrades(ctx context.Context, input *model.BatchGradeInput) ([]model.StudentGrade, error) {
	grades := make([]model.StudentGrade, 0, len(input.Grades))

	for _, g := range input.Grades {
		gradeInput := &model.StudentGradeInput{
			TenantID:        input.TenantID,
			StudentID:       g.StudentID,
			SubjectID:       input.SubjectID,
			SemesterID:      input.SemesterID,
			ScoreNumeric:    g.ScoreNumeric,
			ScorePredicate:  g.ScorePredicate,
			DescriptionHigh: g.DescriptionHigh,
			DescriptionLow:  g.DescriptionLow,
			ComponentScores: g.ComponentScores,
		}

		grade, err := a.SaveGrade(ctx, gradeInput)
		if err != nil {
			return nil, err
		}
		grades = append(grades, *grade)
	}

	return grades, nil
}

func (a *eraporAdapter) GetGradesByStudent(ctx context.Context, studentID, semesterID string) ([]model.StudentGrade, error) {
	var rows []struct {
		ID              string    `db:"id"`
		TenantID        string    `db:"tenant_id"`
		StudentID       string    `db:"student_id"`
		SubjectID       string    `db:"subject_id"`
		SubjectName     string    `db:"subject_name"`
		SemesterID      string    `db:"semester_id"`
		ScoreNumeric    float64   `db:"score_numeric"`
		ScorePredicate  string    `db:"score_predicate"`
		DescriptionHigh string    `db:"description_high"`
		DescriptionLow  string    `db:"description_low"`
		ComponentScores []byte    `db:"component_scores"`
		CreatedAt       time.Time `db:"created_at"`
		UpdatedAt       time.Time `db:"updated_at"`
	}

	err := a.db.From("student_grades").
		Select("student_grades.*", goqu.I("subjects.name").As("subject_name")).
		Join(goqu.T("subjects"), goqu.On(goqu.I("student_grades.subject_id").Eq(goqu.I("subjects.id")))).
		Where(goqu.I("student_grades.student_id").Eq(studentID)).
		Where(goqu.I("student_grades.semester_id").Eq(semesterID)).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, err
	}

	grades := make([]model.StudentGrade, len(rows))
	for i, row := range rows {
		var scores []float64
		if len(row.ComponentScores) > 0 {
			json.Unmarshal(row.ComponentScores, &scores)
		}
		grades[i] = model.StudentGrade{
			ID:              row.ID,
			TenantID:        row.TenantID,
			StudentID:       row.StudentID,
			SubjectID:       row.SubjectID,
			SubjectName:     row.SubjectName,
			SemesterID:      row.SemesterID,
			ScoreNumeric:    row.ScoreNumeric,
			ScorePredicate:  row.ScorePredicate,
			DescriptionHigh: row.DescriptionHigh,
			DescriptionLow:  row.DescriptionLow,
			ComponentScores: scores,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		}
	}

	return grades, nil
}

func (a *eraporAdapter) GetGradesBySubject(ctx context.Context, subjectID, semesterID string) ([]model.StudentGrade, error) {
	var rows []struct {
		ID              string    `db:"id"`
		TenantID        string    `db:"tenant_id"`
		StudentID       string    `db:"student_id"`
		SubjectID       string    `db:"subject_id"`
		SemesterID      string    `db:"semester_id"`
		ScoreNumeric    float64   `db:"score_numeric"`
		ScorePredicate  string    `db:"score_predicate"`
		DescriptionHigh string    `db:"description_high"`
		DescriptionLow  string    `db:"description_low"`
		ComponentScores []byte    `db:"component_scores"`
		CreatedAt       time.Time `db:"created_at"`
		UpdatedAt       time.Time `db:"updated_at"`
	}

	err := a.db.From("student_grades").
		Where(goqu.C("subject_id").Eq(subjectID)).
		Where(goqu.C("semester_id").Eq(semesterID)).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, err
	}

	grades := make([]model.StudentGrade, len(rows))
	for i, row := range rows {
		var scores []float64
		if len(row.ComponentScores) > 0 {
			json.Unmarshal(row.ComponentScores, &scores)
		}
		grades[i] = model.StudentGrade{
			ID:              row.ID,
			TenantID:        row.TenantID,
			StudentID:       row.StudentID,
			SubjectID:       row.SubjectID,
			SemesterID:      row.SemesterID,
			ScoreNumeric:    row.ScoreNumeric,
			ScorePredicate:  row.ScorePredicate,
			DescriptionHigh: row.DescriptionHigh,
			DescriptionLow:  row.DescriptionLow,
			ComponentScores: scores,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		}
	}

	return grades, nil
}

func (a *eraporAdapter) GetStudentRapor(ctx context.Context, studentID, semesterID string) (*model.RaporData, error) {
	// Get student grades
	grades, err := a.GetGradesByStudent(ctx, studentID, semesterID)
	if err != nil {
		return nil, err
	}

	// TODO: Get student info, attendance, extracurricular from their respective tables
	// For now, return basic rapor data
	return &model.RaporData{
		StudentID:       studentID,
		SemesterID:      semesterID,
		Grades:          grades,
		Attendance:      model.AttendanceData{},
		Extracurricular: []model.ExtracurricularData{},
	}, nil
}

func (a *eraporAdapter) GetGradeStats(ctx context.Context, tenantID, semesterID string) (map[string]interface{}, error) {
	// Count total grades
	var totalGrades int64
	_, err := a.db.From("student_grades").
		Select(goqu.COUNT("*")).
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Where(goqu.C("semester_id").Eq(semesterID)).
		ScanValContext(ctx, &totalGrades)
	if err != nil {
		return nil, err
	}

	// Count by predicate
	type predicateCount struct {
		Predicate string `db:"score_predicate"`
		Count     int64  `db:"count"`
	}
	var predicates []predicateCount
	a.db.From("student_grades").
		Select(goqu.C("score_predicate"), goqu.COUNT("*").As("count")).
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Where(goqu.C("semester_id").Eq(semesterID)).
		GroupBy("score_predicate").
		ScanStructsContext(ctx, &predicates)

	predicateMap := make(map[string]int64)
	for _, p := range predicates {
		predicateMap[p.Predicate] = p.Count
	}

	return map[string]interface{}{
		"total_grades":        totalGrades,
		"grades_by_predicate": predicateMap,
		"semester_id":         semesterID,
	}, nil
}
