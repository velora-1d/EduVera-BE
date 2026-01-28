package erapor

import (
	"context"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

// Service adalah domain service untuk E-Rapor
type Service struct {
	db outbound_port.ERaporDatabasePort
}

// NewService membuat instance baru Service
func NewService(db outbound_port.ERaporDatabasePort) *Service {
	return &Service{db: db}
}

// ==========================================
// SUBJECT OPERATIONS
// ==========================================

// CreateSubject membuat mata pelajaran baru dengan konfigurasi grading
func (s *Service) CreateSubject(ctx context.Context, input *model.SubjectInput) (*model.Subject, error) {
	// Apply default grading config if empty
	if len(input.GradingConfig.Components) == 0 {
		if input.Type == model.SubjectTypeFormalK13 {
			input.GradingConfig = model.DefaultK13GradingConfig()
		} else {
			input.GradingConfig = model.DefaultMerdekaGradingConfig()
		}
	}
	return s.db.CreateSubject(ctx, input)
}

// UpdateSubject mengupdate mata pelajaran
func (s *Service) UpdateSubject(ctx context.Context, id string, input *model.SubjectInput) (*model.Subject, error) {
	return s.db.UpdateSubject(ctx, id, input)
}

// GetSubjectByID mengambil mata pelajaran berdasarkan ID
func (s *Service) GetSubjectByID(ctx context.Context, id string) (*model.Subject, error) {
	return s.db.GetSubjectByID(ctx, id)
}

// GetSubjectsByTenant mengambil semua mata pelajaran tenant
func (s *Service) GetSubjectsByTenant(ctx context.Context, tenantID string) ([]model.Subject, error) {
	return s.db.GetSubjectsByTenant(ctx, tenantID)
}

// DeleteSubject menghapus mata pelajaran (soft delete)
func (s *Service) DeleteSubject(ctx context.Context, id string) error {
	return s.db.DeleteSubject(ctx, id)
}

// ==========================================
// GRADE OPERATIONS
// ==========================================

// SaveGrade menyimpan nilai siswa dengan kalkulasi predicate otomatis
func (s *Service) SaveGrade(ctx context.Context, input *model.StudentGradeInput) (*model.StudentGrade, error) {
	// Get subject to apply grading rules
	subject, err := s.db.GetSubjectByID(ctx, input.SubjectID)
	if err != nil {
		return nil, err
	}

	// Auto-calculate predicate if not provided
	if input.ScorePredicate == "" && subject != nil {
		input.ScorePredicate = s.calculatePredicate(input.ScoreNumeric, subject.GradingConfig)
	}

	return s.db.SaveGrade(ctx, input)
}

// BatchSaveGrades menyimpan banyak nilai sekaligus
func (s *Service) BatchSaveGrades(ctx context.Context, input *model.BatchGradeInput) ([]model.StudentGrade, error) {
	// Get subject for predicate calculation
	subject, err := s.db.GetSubjectByID(ctx, input.SubjectID)
	if err != nil {
		return nil, err
	}

	// Auto-calculate predicates
	for i := range input.Grades {
		if input.Grades[i].ScorePredicate == "" && subject != nil {
			input.Grades[i].ScorePredicate = s.calculatePredicate(input.Grades[i].ScoreNumeric, subject.GradingConfig)
		}
	}

	return s.db.BatchSaveGrades(ctx, input)
}

// GetGradesByStudent mengambil semua nilai siswa di semester tertentu
func (s *Service) GetGradesByStudent(ctx context.Context, studentID, semesterID string) ([]model.StudentGrade, error) {
	return s.db.GetGradesByStudent(ctx, studentID, semesterID)
}

// GetGradesBySubject mengambil semua nilai untuk mata pelajaran tertentu
func (s *Service) GetGradesBySubject(ctx context.Context, subjectID, semesterID string) ([]model.StudentGrade, error) {
	return s.db.GetGradesBySubject(ctx, subjectID, semesterID)
}

// GetStudentRapor mengambil data rapor lengkap untuk siswa
func (s *Service) GetStudentRapor(ctx context.Context, studentID, semesterID string) (*model.RaporData, error) {
	return s.db.GetStudentRapor(ctx, studentID, semesterID)
}

// GetGradeStats mengambil statistik nilai
func (s *Service) GetGradeStats(ctx context.Context, tenantID, semesterID string) (map[string]interface{}, error) {
	return s.db.GetGradeStats(ctx, tenantID, semesterID)
}

// ==========================================
// UTILITY FUNCTIONS
// ==========================================

// calculatePredicate menghitung predikat berdasarkan nilai dan config
func (s *Service) calculatePredicate(score float64, config model.GradingConfig) string {
	scoreInt := int(score)
	for _, rule := range config.PredicateRules {
		if scoreInt >= rule.MinScore && scoreInt <= rule.MaxScore {
			return rule.Predicate
		}
	}
	// Default fallback
	if scoreInt >= 90 {
		return "A"
	} else if scoreInt >= 80 {
		return "B"
	} else if scoreInt >= 70 {
		return "C"
	}
	return "D"
}

// CalculateFinalScore menghitung nilai akhir dari component scores
func (s *Service) CalculateFinalScore(componentScores []float64, config model.GradingConfig) float64 {
	if len(componentScores) == 0 || len(config.Components) == 0 {
		return 0
	}

	var total float64
	var totalWeight int

	for i, component := range config.Components {
		if i < len(componentScores) {
			total += componentScores[i] * float64(component.Weight)
			totalWeight += component.Weight
		}
	}

	if totalWeight == 0 {
		return 0
	}

	return total / float64(totalWeight)
}

// GenerateDescription menghasilkan deskripsi capaian kompetensi
func (s *Service) GenerateDescription(score float64, subjectName string) (high, low string) {
	if score >= 90 {
		high = "Ananda menunjukkan penguasaan yang sangat baik pada materi " + subjectName
		low = "Tetap pertahankan dan tingkatkan terus prestasi belajar"
	} else if score >= 80 {
		high = "Ananda menunjukkan pemahaman yang baik terhadap konsep-konsep " + subjectName
		low = "Perlu latihan lebih untuk mencapai hasil yang maksimal"
	} else if score >= 70 {
		high = "Ananda memahami konsep dasar " + subjectName
		low = "Perlu bimbingan tambahan untuk penguatan materi"
	} else {
		high = "Ananda sudah berusaha dalam mempelajari " + subjectName
		low = "Perlu pendampingan intensif untuk peningkatan pemahaman"
	}
	return high, low
}
